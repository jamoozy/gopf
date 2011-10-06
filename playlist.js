var playlist = (function() {
  var dir = "data/";
  var req = new XMLHttpRequest();

  function notify(str) {
    document.getElementById("notification").innerHTML = str;
  }

  function reqPlaylist(elem) {
    var url = document.location.pathname + "list.php?playlist=" +
        elem.innerHTML.replace(/^\s+|\s+$/g,"");
    req.open("GET", url, true);
    req.send();
  }

  function loadPlaylist(req) {
    var media = req.responseText.replace(/\.\.\//g, dir + "/").split("\n");
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
    for (i = 0; i < media.length; i++) {
      if (media[i].length > 0) {
        var media_first = media[i].lastIndexOf("/") + 1;
        var media_length = media[i].lastIndexOf(".") - media_first;
        var name = media[i].substr(media_first, media_length);

        mediaTag = document.createElement("li");
        mediaTag.setAttribute("class", "media");
        mediaTag.setAttribute("path", escape(media[i]));
        mediaTag.setAttribute("onclick", "media.onclick(this)");
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
    }
  };
})();
