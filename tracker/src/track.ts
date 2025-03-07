export const ENDPOINT_URL: string = "http://localhost:3000/track"

export interface ITrackingData {
  url?: string;
  visitorId: string;
  trackingId: string;
  referrer: string | null;
  ua: string;
  details?: Record<string, any>;
}

export interface EventPayload {
  tracking: ITrackingData;
  type: string;
}

export class Analytics {
  private visitorId: string | null = null;
  private trackingId: string;
  private referrer: string | null;


  constructor(trackingId: string) {
    if (!trackingId) {
      throw new Error("Tracking ID is required");
    }
    this.trackingId = trackingId;
    this.referrer = this.getReferrer();
  }

  public async track(type: string = 'pageview', details?: Record<string, any>) {
    if (!this.visitorId) {
      this.visitorId = await this.generateDailyVisitorHash();
    }
    const trackingData: ITrackingData = {
      visitorId: this.visitorId,
      trackingId: this.trackingId,
      referrer: this.referrer,
      ua: navigator.userAgent,
    };
    if (type === "pageview") {
      trackingData.url = window.location.href;
    }
    if (details) {
      trackingData.details = details;
    }
    const payload: EventPayload = {
      tracking: trackingData,
      type: type
    };
    this.sendData(payload);
  }

  public trackSubsequentPages() {
    const originalPushState = window.history.pushState;
    const originalReplaceState = window.history.replaceState;
    window.history.pushState = (...args) => {
      originalPushState.apply(window.history, args);
      this.track();
    };
    window.history.replaceState = (...args) => {
      originalReplaceState.apply(window.history, args);
      this.track();
    };
    window.addEventListener('popstate', () => {
      this.track();
    });
  }

  private async generateDailyVisitorHash(): Promise<string> {
    const uaHash = await this.createHash(navigator.userAgent);
    const components = [uaHash, new Date().toLocaleDateString()].join("|");
    return this.createHash(components);
  }

  private async createHash(input: string): Promise<string> {
    const encoder = new TextEncoder();
    const data = encoder.encode(input);
    const hashBuffer = await crypto.subtle.digest("SHA-256", data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, "0")).join("");
    return hashHex;
  }

  private getReferrer() {
    const ref = document.referrer;
    return ref && new URL(ref).host !== window.location.hostname ? ref : null;
  }

  private sendData(payload: EventPayload) {
    const s = JSON.stringify(payload);
    const url = `${ENDPOINT_URL}?data=${btoa(s)}`;
    const img = new Image();
    img.onerror = () => {
      console.error("Failed to send analytics data");
      return;
    };
    img.src = url;
  }
}

// Initialize analytics
(async (w: Window & { _analytics?: Analytics }, d: Document) => {
  const script = d.currentScript as HTMLScriptElement;
  const trackingId = script?.getAttribute("tracking-id");
  if (!trackingId) {
    console.error("Analytics: No tracking ID provided");
    return;
  }
  const analytics = new Analytics(trackingId);
  await analytics.track();
  analytics.trackSubsequentPages();
  w._analytics = analytics;
})(window, document);
