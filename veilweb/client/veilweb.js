
/*
Meteor.subscribe('posts', function() {
  console.log(Posts.find().count());
});
Alternatively, you can call ready on the subscription handle:

var handle = Meteor.subscribe('posts');

Tracker.autorun(function() {
  if (handle.ready())
    console.log(Posts.find().count());
});

*/

var fragments_loaded = false;

Meteor.subscribe("fragments", function(){
    fragments_loaded = true
});

Tracker.autorun(function () {
  //if(frag_handle.ready()) {
    console.log("Fragments ready: ", Fragments.find().count())
  //}
});

Fragments.find().observe({
    added: function (fragment) {
        console.log("Added: ", fragment.text);
        addBranchFromStringArray(fragment.text.split(","));
    }
});


