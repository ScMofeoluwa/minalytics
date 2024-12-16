import { createHash } from "crypto";

interface ITrackingData {
  url?: string;
  visitorId: string,
  referrer: string | null,
  ua: string,
  details?: Record<string, any>
}

interface EventPayload {
  tracking: ITrackingData,
  type: string
}

class Analytics {
  private visitorId: string;
  private referrer: string | null;

  constructor(){
    this.visitorId = this.generateDailyVisitorHash()
    this.referrer = this.getReferrer()
  }

  public track(type: string = 'pageview', details?: Record<string, any>) {
    const trackingData: ITrackingData = {
      visitorId: this.visitorId,
      referrer: this.referrer,
      ua: navigator.userAgent,
    }
    
    if (type === "payload"){
      trackingData.url = window.location.href
    }

    if (details) {
      trackingData.details = details
    }

    const payload: EventPayload = {
      tracking: trackingData,
      type: type
    }

    this.sendData(payload)
  }

  public trackSubsequentPages(){
    const originalPushState = window.history.pushState
    const originalReplaceState = window.history.replaceState

    window.history.pushState = (...args) => {
      originalPushState.apply(this, args)
      this.track()
    }

    window.history.replaceState = (...args) => {
      originalReplaceState.apply(this, args)
      this.track()
    }

    window.addEventListener('popstate', () => {
      this.track()
    })

    // window.addEventListener('hashchange', () => {
    //   this.track(document.location.hash)
    // })
  }

  private generateDailyVisitorHash(): string{
    const components = [
      navigator.userAgent,
      new Date().toDateString()
    ].join("|")

    return this.createHash(components)
  }

  private createHash(input: string): string{
    return createHash('sha1').update(input).digest('hex')
  }

  private getReferrer() {
    const ref = document.referrer
    return ref && new URL(ref).host !== window.location.hostname ? ref : null
  }

  private sendData(payload: EventPayload){
    const s = JSON.stringify(payload)
    const url = `http://localhost:3000/track?data=${btoa(s)}`

    const img = new Image()
    img.src = url;
  }
}

((w,d) => {
  const analytics = new Analytics()
  analytics.track()
  analytics.trackSubsequentPages()

  w._analytics = analytics
})
