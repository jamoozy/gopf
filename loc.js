// Copyright 2012 Andrew "Jamoozy" Correa
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// Foobar is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License
// for more details.
//
// You should have received a copy of the GNU General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.

var loc = (function() {
  function deletePage() {
    var body = document.getElementsByTagName("body")[0];
    body.innerHTML = "<div style=\"text-align:center;margin-top:40px\"><h1> NO SOUP FOR YOU </h1></div>"; }

  function errorHandler() {
    window.console.log("Therer was an error.");
    deletePage();
  }

  function positionHandler(pos) {
    logPos(loc.pos = pos);
  }

  return {
    pos : {},

    init : function(e) {
      if (!navigator.geolocation) {
        deletePage();
      }

      navigator.geolocation.getCurrentPosition(positionHandler, errorHandler);
    }
  };
})();

window.addEventListener("load", loc.init, true);
