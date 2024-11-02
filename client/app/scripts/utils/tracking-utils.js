import debug from 'debug';

const log = debug('service:tracking');

// Track segment events only if Scope is running inside of nholuongut Cloud.
export function trackAnalyticsEvent(name, props) {
  if (window.analytics && process.env.nholuongut_CLOUD) {
    window.analytics.track(name, props);
  } else {
    log('trackAnalyticsEvent', name, props);
  }
}
