import { createHash } from "crypto";

interface ITrackingData {
  url: string;
  visitorId: string,
  referrer: string | null
  deviceType: string,
  browserName: string,
  os: string
}

interface TrackPayload {
  tracking: ITrackingData,
  type: string
}

class Analytics {
  private visitorId: string;

  constructor(){
    this.visitorId = this.generateDailyVisitorHash()
  }

  public track(type: string = 'pageview') {
    const ref = document.referrer
    const externalReferrer = ref && new URL(ref).host !== window.location.hostname ? ref : null

    const trackingData: ITrackingData = {
      url: window.location.href,
      visitorId: this.visitorId,
      referrer: externalReferrer,
      deviceType: this.getDeviceType(),
      browserName: this.getBrowserName(),
      os: this.getOS()
    }

    const payload: TrackPayload = {
      tracking: trackingData,
      type: type
    }

    this.sendData(payload)
  }

  private generateDailyVisitorHash(): string{
    const components = [
      navigator.userAgent,
      new Date().toDateString()
    ].join("|")

    return this.createHash(components)
  }

  private getDeviceType(): string{
    const isTablet = window.matchMedia('(min-width: 768px) and (max-width: 1024px)').matches
    const isMobile = window.matchMedia('(max-width: 767px)').matches || ('maxTouchPoints' in navigator && navigator.maxTouchPoints > 0)

    if (isTablet) return 'tablet'
    if (isMobile) return 'mobile'
    return 'desktop'
  }

  private getBrowserName(): string{
    const userAgentData = (navigator as any).userAgentData
    const ua = navigator.userAgent

    if (userAgentData?.brands) {
      const brand = userAgentData.brands.find(
        (b: {brand: string}) => !["Chromium", "Not-A.Brand"].includes(b.brand)
      )
      if (brand) return brand.brand
    }

    if (ua.includes("Firefox/")){
      return 'Firefox'
    } else if ((window as any).chrome) {
      if (ua.includes("Edg/")) return "Edge"
      if (ua.includes("OPR/")) return "Opera"
      return "Chrome"
    } else if (ua.includes("Safari/") && !ua.includes("Chrome/")){
      return "Safari"
    }
    return "Unknown"
  }

  private getOS(): string {
    const userAgentData = (navigator as any).userAgentData

    if (userAgentData?.platform){
      return userAgentData.platform
    }

    const platform = navigator.platform

    if (platform.includes("Win")) return "Window" 
    if (platform.includes("Mac")) return "MacOS" 
    if (/iPhone|iPad|iPod/.test(platform)) return "iOS"
    if (platform.includes("Linux")) {
      return navigator.userAgent.includes("Android") ? "Android": "Linux"
    }
    
    return "Unknown"
  }

  private createHash(input: string): string{
    return createHash('sha1').update(input).digest('hex')
  }

  private sendData(payload: TrackPayload){
    const s = JSON.stringify(payload)
    const url = `http://localhost:3000/track?data=${btoa(s)}`

    const img = new Image()
    img.src = url;
  }
}
