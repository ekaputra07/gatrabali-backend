const firesub = require("firesub");
const topic = "firestore_events";

/**
 * Below are all Firestore triggered events handler, instead of process them directly
 * we publish them to a PubSub topic, and our worker service will then handle and process them.
 *
 * Please note that we're only using single topic here since the event won't be that much.
 * For better scalability please consider using separate topic.
 */
exports.entryOnCreate = firesub.FirestoreOnCreate("/entries/{entryId}", topic, {
  collection: "entries"
});

exports.kriminalOnCreate = firesub.FirestoreOnCreate(
  "/kriminal/{entryId}",
  topic,
  {
    collection: "kriminal"
  }
);

exports.baliunitedOnCreate = firesub.FirestoreOnCreate(
  "/baliunited/{entryId}",
  topic,
  {
    collection: "baliunited"
  }
);

exports.balebengongOnCreate = firesub.FirestoreOnCreate(
  "/balebengong/{entryId}",
  topic,
  {
    collection: "balebengong"
  }
);

exports.entryResponseOnWrite = firesub.FirestoreOnWrite(
  "/entry_responses/{entryId}",
  topic,
  {
    collection: "entry_responses"
  }
);
