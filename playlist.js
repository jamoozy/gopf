// Copyright (C) 2011-2015 Andrew "Jamoozy" Sabisch
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under the
// terms of the GNU Affero General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more
// details.
//
// You should have received a copy of the GNU Affero General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.

var playlist = (function() {
  var dir = "data/",
      req = new XMLHttpRequest(),
      callback = false;

  function notify(str) {
    $("#notification").html(str);
  }

  // Requests the contents of the playlist from the server.
  function reqPlaylist(elem) {
    window.console.log("reqPlaylist(", elem, ")");
    var html = elem.html();
    html.replace(/^\s+|\s+$/g,"");
    var url = document.location.pathname + "list?playlist=" + html;
    req.open("GET", url, true);
    req.send();
  }

  function loadPlaylist(req) {
    var path = req.responseText.replace(/\.\.\//g, dir).split("\n");
    var queue = $("#media");
    var mediaTag, i;

    // Remove sources.
    $("#player").removeAttr("src");

    // Remove all current children.
    queue.children().each(function(i,e) {
      e.remove();
    });

    // Add the new children.
    for (i = 0; i < path.length; i++) {
      if (path[i].trim().length <= 0 || path[i][0] == "#") {
        continue;
      }
      var media_first = path[i].lastIndexOf("/") + 1;
      var media_length = path[i].lastIndexOf(".") - media_first;
      var name = path[i].substr(media_first, media_length);

      mediaTag = $("<li>");
      mediaTag.attr("class", "media");
      mediaTag.attr("path", path[i]);
      mediaTag.click(function(e) {
        media.onclick(this);
      });
      mediaTag.html(name);

      queue.append(mediaTag);
    }

    // Ensure the media queue isn't overlapping things.
    $("#media-container").width = (window.innerWidth - $("#playlist-container").width) / 2;

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
      $(".unselected").click(function(event) {
          window.console.log("A playlist was clicked: " + this);
          $("#selected").removeAttr("id");
          $(this).attr("id", "selected");
          playlist.onclick(this);
      });
    },

    // Register that a playlist was clicked (to be loaded).
    //        elem: The clicked element.
    //   swapAfter: Whether the list the user is controlling should be
    //              swapped after it is loaded.
    //          cb: (optional) Callback after the list is loaded.
    onclick : function(elem, swapAfter, cb) {
      console.dir('elem is:', elem);
      if (!!elem.target) {
        console.log('elem has a target');
        elem = elem.target;
      }
      elem = $(elem);
      if (elem.attr("class") === "selected") {
        return;
      }

      if (swapAfter) {
        playlist.swapAfter = true;
      }

      var selected = $(".selected");
      for (var i = 0; i < selected.size(); i++) {
        $(selected.get(i)).attr("class", "unselected");
      }
      elem.attr("class", "selected");
      reqPlaylist(elem);
      setCallback(cb);
    }
  };
})();

$(window).load(playlist.init);
