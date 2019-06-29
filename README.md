#  Gatra Bali Backend
Backend of Gatra Bali app including:
- _common_: a packaged that shared between cloud functions and api project
- _cloud-functions_: a collection of Cloud Functions
- _api_: a Google Appengine app for the app REST api.

# How it works
<img src="https://raw.githubusercontent.com/apps4bali/gatrabali-backend/master/arch.jpg"/>

1. [Miniflux](https://github.com/apps4bali/miniflux) (an opensource Feed reader written in Go) will periodically check for new articles on a given feeds. Since its quite full featured we're able to add new feed sources, categorize the feed, manage users, etc. Miniflux works independently know nothing about the other parts, store its own data to its own database (PostgreSQL).

1. The app didn't talk directly to Miniflux Api even though its possible but I want Firestore as the data storage instead. So I need a way to transfer data from Miniflux to Firestore. Here I utilise the Google PubSub to trigger the **SyncData** cloud function, whenever an Article is added/updated, Feed is added/updated/deleted, Category is added/updated/deteled Miniflux will publish a message to a Topic and it will trigger the Cloud Function to running. 

1. When the PubSub triggered Function is running, based on the message it received it will make a HTTP request back to the Miniflux REST api to get the Article, Feed or Category object and store them to Firestore.

1. Data that received from Miniflux is stored in Firestore on separate collections, eg. categories, feeds, entries.

1. On the app we allow users to subscribe to a category and receive push notification when new article published in that category, here the **NotifyCategorySubscribers** function is triggered whenever new article/entry written to Firestore. This function doesn't do the actual sending of push notification but only doing a check on who are currently subscribes to a category and publish a push message to *PushNotification* pub/sub topic for each of them (subscribers).

1. **SendPushNotification** cloud function is sitting there alone listen for new message pushed to *PushNotification* topic, this function doesn't know anything about our data/Firestore, it only do one thing which is sending push notification to users. 

1. Because most of the data will be public (user doesn't need to register to use the app except for bookmarks) and to reduce the Firestore reads, REST API based on Appengine is used to serve the requests from the app, some of the advantages of this is that we can implement caching and auto-scaling.

1. For the bookmarks, apps need to read and write directly to Firestore, users are only allowed to read and write to their own bookmarks collection by using Firestore rules. User's Cloud Messaging tokens are also written directly to Firestore by the app.
