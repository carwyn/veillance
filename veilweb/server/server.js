
Meteor.publish("fragments", function () { 
    return Fragments.find();
});

Meteor.startup(function () {

  Accounts.loginServiceConfiguration.remove({
    service : 'twitter'
  });

  Accounts.loginServiceConfiguration.insert({
    service     : 'twitter',
    consumerKey : 'Qb5PBQCTojNVcyweXSxWCoHOx',
    secret      : 'hoIGjM8Me3HXky7PQbB8sFaE0GIOMnfgiMUGKPexCWVVg2JnoK'
  });
});

