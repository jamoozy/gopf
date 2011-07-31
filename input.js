var input = (function() {
  // Extra amount to shrink media container by, so that media/playlists
  // don't overlap one onther.
  var DIVIDER_WIDTH = 40;

  // Whether the selector is in the media list.  False means it's in the
  // playlists list.
  var viewMedia = false;

  function notify(str) {
    document.getElementById("notification").innerHTML = str;
  }

  // Select
  function select(elem) {
    var sel = document.getElementById("selected");
    if (sel !== null && sel !== undefined) {
      sel.removeAttribute("id");
    }
    if (elem.setAttribute !== undefined) {
      elem.setAttribute("id", "selected");
    } else {
      getElems()[0].setAttribute("id", "selected");
    }
  }

  // Gets the elements from the currently-selected list.
  function getElems() {
    return getList().childNodes;
  }

  // Gets the currently-selected list.
  function getList() {
    if (viewMedia) {
      return document.getElementById("songs");
    } else {
      return document.getElementById("playlists");
    }
  }

  // Gets the index of an element in the list.
  function getElemIndex(list, elem) {
    if (elem === null || elem === undefined) { return 0; }
    for (var i = 0; i < list.length; i++) {
      if (list[i] === elem) {
        return i;
      }
    }
    return 0;
  }

  // Ensures the selected element is visible, either by centering the
  // list's view on the element, or scrolling the list all the way up or
  // down so the element is in view.
  function ensureSelectedVisible(i) {
    var sel = document.getElementById("selected");
    var list = getList();
    // "-1" due to "border-bottom: -1px"
    list.scrollTop = i * (sel.offsetHeight - 1) - list.offsetHeight / 2;
  }

  // Selects the previous element in the list.
  function prev(dec) {
    var sel = document.getElementById("selected");
    var elems = getElems();
    var i = getElemIndex(elems, sel);

    if (dec === undefined) {
      dec = 1;
    }

    if (!isNaN(i)) {
      var val = i - dec;
      if (val < 0) {
        val = 0;
      }
      select(elems[val]);
      ensureSelectedVisible(val);
    }
  }

  // Selects the next element in the list.
  function next(inc) {
    var sel = document.getElementById("selected");
    var elems = getElems();
    var i = getElemIndex(elems, sel);

    if (inc === undefined) {
      inc = 1;
    }

    if (!isNaN(i)) {
      var val = i + inc;
      if (val >= elems.length) {
        val = elems.length - 1;
      }
      select(elems[val]);
      ensureSelectedVisible(val);
    }
  }

  // Key handler.
  function onkey(e) {
    switch(e.keyCode) {
      case 33:  // PageUp
        prev(10);
        e.stopPropagation();
        break;
      case 38:  // up
      case 75:  // k
        prev(1);
        e.stopPropagation();
        break;
      case 34:  // PageDown
        next(10);
        e.stopPropagation();
        break;
      case 40:  // down
      case 74:  // j
        next(1);
        e.stopPropagation();
        break;
      case 192: // back-tick ("`")
      case 72:  // h
        input.swap();
        e.stopPropagation();
        break;

      case 80:  // P
      case 37:  // left
        media.prev();
        break;
      case 78:  // N
      case 39:  // right
        media.next();
        break;

      case 83:  // S
        document.getElementById("shuf").click();
        break;
      case 76:  // L
        document.getElementById("loop").click();
        break;

      case 77:  // M
        var player = document.getElementById("player");
        player.muted = !player.muted;
        break;
      case 32:  // space
        var player = document.getElementById("player");
        if (player.paused) {
          player.play();
        } else {
          player.pause();
        }
        break;
      case 13:  // enter
        var sel = document.getElementById("selected");
        if (sel !== null && sel !== undefined) {
          if (viewMedia) {
            media.onclick(sel);
          } else {
            playlist.onclick(sel, true);
          }
        }
        break;
    }
  }

  // Window resize handler.
  function adjustSize() {
    var song_cont = document.getElementById("song-container");
    var songs = document.getElementById("songs");
    var song_head = document.getElementById("song-header");
    var playlist_cont = document.getElementById("playlist-container");
    var playlists = document.getElementById("playlists");
    var playlist_head = document.getElementById("playlist-header");
    var media = document.getElementById("media-container");
    var footer = document.getElementById("footer");

    var width = window.innerWidth - playlist_cont.offsetLeft -
                playlist_cont.offsetWidth - DIVIDER_WIDTH;
    var height = footer.offsetTop - media.offsetHeight -
                 media.offsetTop - playlist_cont.offsetLeft;

    //notify("resized to " + width);
    song_cont.style.maxWidth = width + "px";
    song_cont.style.maxHeight = height + "px";
    songs.style.maxHeight = height - song_head.offsetHeight + "px";
    playlist_cont.style.maxHeight = height + "px";
    playlists.style.maxHeight = height - playlist_head.offsetHeight + "px";
  }

  return {
    // Initializes the input module by registering event listeners.
    init : function() {
      var sel = document.getElementById("selected");
      if (sel === null) {
        var playlists = document.getElementById("playlists").childNodes;
        if (playlists.length > 0) {
          playlists[0].setAttribute("id", "selected");
        }
      }

      window.onresize = function(e) { adjustSize(); } 
      window.onkeydown = onkey;
      adjustSize();
    },

    // Swaps the list the user is moving through.
    swap : function() {
      viewMedia = !viewMedia;
      if (viewMedia) {
        if (document.getElementsByClassName("dummy").length > 0) {
          viewMedia = false;
        } else {
          select(document.getElementById("songs").childNodes[0]);
        }
      } else {
        select(document.getElementById("playlists").childNodes[0]);
      }
    }
  };
})();

window.addEventListener("load", input.init, true);
