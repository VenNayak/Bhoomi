//SPDX-License-Identifier: Apache-2.0

var bhoomi = require('./controller.js');

module.exports = function(app){

  app.get('/get_landRecord/:id', function(req, res){
    bhoomi.get_landRecord(req, res); // fetch Land Records via MOJANI/KAVERI or finance application
  });
  app.get('/create_landRecord/:landRecord', function(req, res){
    bhoomi.create_landRecord(req, res); // create new land record via MOJANI
  });
  app.get('/transfer_ownership/:owner', function(req, res){
    bhoomi.transfer_ownership(req, res); //transfer ownership of land records via KAVERI
  });
  app.get('/allot_landRecord/:owner', function(req, res){
    bhoomi.allot_landRecord(req, res);   //allot a new land record to first owner via MOJANI
  });
}
