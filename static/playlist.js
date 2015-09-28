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
  var dir = "data/";
  var callback = false;

  function notify(str) {
    $("#notification").html(str);
  }

  // Requests the contents of the playlist from the server.
  function reqPlaylist(elem) {
    window.console.log("reqPlaylist(...):");
    window.console.log(elem);

    $.get("playlist/" + elem.html(), null, function(data, textStatus, jqXHR) {
      window.console.log("Got data:");
      window.console.log(data);

      var queue = $("#media"),
          mediaTag,
          i;

      // Remove sources.
      $("#player").removeAttr("src");

      // Remove all current children.
      queue.children().each(function(i,e) {
        e.remove();
      });

      // Add the new children.
      for (i = 0; i < data.Files.length; i++) {
        if (data.Files[i].trim().length <= 0 || data.Files[i][0] == "#") {
          continue;
        }
        var media_first = data.Files[i].lastIndexOf("/") + 1;
        var media_length = data.Files[i].lastIndexOf(".") - media_first;
        var name = data.Files[i].substr(media_first, media_length);

        mediaTag = $("<li>");
        mediaTag.attr("class", "media");
        mediaTag.attr("path", data.Files[i]);
        mediaTag.click(function(e) {
          media.onclick(this);
        });
        mediaTag.html(name);

        queue.append(mediaTag);
      }

      // Ensure the media queue isn't overlapping things.
      $("#media-container").width = (
          window.innerWidth - $("#playlist-container").width) / 2;

      // If we need to swap the selection, do it.
      if (playlist.swapAfter) {
        input.swap();
        playlist.swapAfter = false;
      }
    }, "json");
  }

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
