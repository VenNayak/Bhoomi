//SPDX-License-Identifier: Apache-2.0

var bhoomi = require('./controller.js');

module.exports = function(app){

  app.get('/get_landRecord/:id', function(req, res){
    bhoomi.get_landRecord(req, res);
  });
  app.get('/create_landRecord/:landRecord', function(req, res){
    bhoomi.create_landRecord(req, res);
  });
  app.get('/transfer_ownership/:owner', function(req, res){
    bhoomi.transfer_ownership(req, res);
  });
  app.get('/init_ownership/:owner', function(req, res){
    bhoomi.init_ownership(req, res);
  });
}
