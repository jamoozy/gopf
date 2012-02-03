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

var playlist = (function() {
  var dir = "data/";
  var req = new XMLHttpRequest();
  var callback = false;

  function notify(str) {
    document.getElementById("notification").innerHTML = str;
  }

  // Requests the contents of the playlist from the server.
  function reqPlaylist(elem) {
    var url = document.location.pathname + "list.php?playlist=" +
        elem.innerHTML.replace(/^\s+|\s+$/g,"");
    req.open("GET", url, true);
    req.send();
  }

  function loadPlaylist(req) {
    var path = req.responseText.replace(/\.\.\//g, dir).split("\n");
    var queue = document.getElementById("media");
    var player = document.getElementById("player");
    var mediaTag, i;

    // Remove sources.
    player.removeAttribute("src");

    // Remove all current children.
    while (queue.childNodes.length > 0) {
      queue.removeChild(queue.firstChild);
    }

    // Add the new children.
    for (i = 0; i < path.length; i++) {
      if (path[i].length > 0) {
        var media_first = path[i].lastIndexOf("/") + 1;
        var media_length = path[i].lastIndexOf(".") - media_first;
        var name = path[i].substr(media_first, media_length);

        mediaTag = document.createElement("li");
        mediaTag.setAttribute("class", "media");
        mediaTag.setAttribute("path", escape(path[i]));
        mediaTag.addEventListener("click", function(event) {
            media.onclick(this);
        }, true);
        mediaTag.innerHTML = name;

        queue.appendChild(mediaTag);
      }
    }

    // Ensure the media queue isn't overlapping things.
    document.getElementById("media-container").width = (window.innerWidth - document.getElementById("playlist-container").width) / 2;

    // If we need to swap the selection, do it.
    if (playlist.swapAfter) {
      input.swap();
      playlist.swapAfter = false;
    }
  }

  // Set the callback for the request.
  req.onreadystatechange = function() {
    switch (req.readyState) {
      case 0: break;
      case 1: break;
      case 2: break;
      case 3: break;
      case 4:
        if (req.status === 200) {
          loadPlaylist(req);
        } else {
          // error?
        }
        break;
      default:
        notify("Not sure what happened ... (default) ... error?");
    }

    if (callback) {
      callback(req);
    }
  };

  // Sets the callback after the list is loaded.
  function setCallback(cb) {
    console.log("set cb");
    callback = cb;
  };

  return {
    swapAfter : false,

    init : function() {
      // Initialize the playlists' "onclick" events.
      var unselected = document.getElementsByClassName("unselected");
      for (var i = 0; i < unselected.length; i++) {
        unselected[i].addEventListener("click", function(event) {
            window.console.log("A playlist was clicked: " + this);
            playlist.onclick(this);
        }, true);
      }
    },

    // Register that a playlist was clicked (to be loaded).
    //        elem: The clicked element.
    //   swapAfter: Whether the list the user is controlling should be
    //              swapped after it is loaded.
    //          cb: (optional) Callback after the list is loaded.
    onclick : function(elem, swapAfter, cb) {
      if (elem.getAttribute("class") === "selected") {
        return;
      }

      if (swapAfter === true) {
        playlist.swapAfter = true;
      }

      var selected = document.getElementsByClassName("selected");
      for (var i = 0; i < selected.length; i++) {
        selected[i].setAttribute("class", "unselected");
      }
      elem.setAttribute("class", "selected");
      reqPlaylist(elem);
      setCallback(cb);
    }
  };
})();

window.addEventListener("load", playlist.init, true);
