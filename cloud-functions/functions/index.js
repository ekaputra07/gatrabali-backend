const firesub = require("firesub");
const topic = "FirestoreEvents";

/**
 * Below are all Firestore triggered events handler, instead of process them directly
 * we publish them to a PubSub topic, and our server will then handle and process them.
 *
 * Please note that we're only using single topic here since the event won't be that much.
 * For better scalability please consider using separate topic.
 */
exports.entryOnCreate = firesub.FirestoreOnCreate("/entries/{entryId}", topic, {
  type: "entries"
});

exports.kriminalOnCreate = firesub.FirestoreOnCreate(
  "/kriminal/{entryId}",
  topic,
  {
    type: "entries"
  }
);

exports.baliunitedOnCreate = firesub.FirestoreOnCreate(
  "/baliunited/{entryId}",
  topic,
  {
    type: "entries"
  }
);

exports.balebengongOnCreate = firesub.FirestoreOnCreate(
  "/balebengong/{entryId}",
  topic,
  {
    type: "entries"
  }
);

exports.entryResponseOnWrite = firesub.FirestoreOnWrite(
  "/entry_responses/{entryId}",
  topic,
  {
    type: "responses"
  }
);
