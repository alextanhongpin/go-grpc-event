db.createUser({
  user: 'user',
  pwd: 'password',
  roles: [ {
    role: 'readWrite',
    db: 'engineersmy'
  }, {
    role: 'dbOwner',
    db: 'engineersmy'
  } ]
}, {
  w: 'majority',
  wtimeout: 5000
});
db.createCollection('events');
