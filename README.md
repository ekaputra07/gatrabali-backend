#  BaliFeed Backend
Backend of BaliFeed app including:
- **[Cloud Functions](https://cloud.google.com/functions)**, responsible for publishing Firestore events to [PubSub](https://cloud.google.com/pubsub).
- **Firebase Hosting**, hosting configuration that rewrite traffic into Cloud Run instance, as well as a CDN.
- **Server**, the core service. Act as API server as well as async workers runs on [Cloud Run](https://cloud.google.com/run).

# How it works
<img src="https://raw.githubusercontent.com/apps4bali/gatrabali-backend/master/backend.png"/>

Check here for the interactive version of the image above: https://whimsical.com/5H3gLFefBuWGRLkG28X3bs
