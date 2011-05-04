var playlist = (function() {
  var root = "/~media/audio/";
  var dir = "data/";
  var req = new XMLHttpRequest();

  function notify(str) {
    document.getElementById("notification").innerHTML = str;
  }

  function reqPlaylist(elem) {
    var url = root + "list.php?playlist=" + elem.innerHTML.replace(/^\s+|\s+$/g,"");
    req.open("GET", url, true);
    req.send();
  }

  function loadPlaylist(req) {
    var songs = req.responseText.replace(/\.\.\//g, dir + "/").split("\n");
    var queue = document.getElementById("songs");
    var player = document.getElementById("player");
    var songTag, i;

    // Remove sources.
    player.removeAttribute("src");

    // Remove all current children.
    while (queue.childNodes.length > 0) {
      queue.removeChild(queue.firstChild);
    }

    // Add the new children.
    for (i = 0; i < songs.length; i++) {
      if (songs[i].length > 0) {
        var song_lio = songs[i].lastIndexOf("/");
        var name = songs[i].substr(1 + song_lio, songs[i].length - song_lio - 5);

        songTag = document.createElement("li");
        songTag.setAttribute("class", "song");
        songTag.setAttribute("path", songs[i]);
        songTag.setAttribute("onclick", "audio.onclick(this)");
        songTag.innerHTML = name;

        queue.appendChild(songTag);
      }
    }

    // ensure the song queue isn't overlapping things.
    document.getElementById("song-container").width = (window.innerWidth - document.getElementById("playlist-container").width) / 2;

    // If we need to swap the selection, do it.
    if (playlist.swapAfter) {
      input.swap();
      playlist.swapAfter = false;
    }
  }


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
  };

  return {
    swapAfter : false,

    onclick : function(elem, swapAfter) {
      if (elem.getAttribute("class") === "selected") { return; }

      if (swapAfter === true) { playlist.swapAfter = true; }

      var selected = document.getElementsByClassName("selected");
      for (var i = 0; i < selected.length; i++) {
        selected[i].setAttribute("class", "unselected");
      }
      elem.setAttribute("class", "selected");
      reqPlaylist(elem);
    }
  };
})();
