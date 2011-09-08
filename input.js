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
      return document.getElementById("media");
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

      case 219:  // [
        document.getElementById("player").playbackRate *= 0.5;
        break;
      case 221:  // ]
        document.getElementById("player").playbackRate *= 2;
        break;
      case 8:  // backspace
        document.getElementById("player").playbackRate = 1.0;
        break;

      case 49:  // 1
        document.getElementById("player").width = 200;
        adjustSize();
        break;
      case 50:  // 2
        document.getElementById("player").width = 400;
        adjustSize();
        break;
      case 51:  // 3
        document.getElementById("player").width = 600;
        adjustSize();
        break;
      case 52:  // 4
        document.getElementById("player").width = 800;
        adjustSize();
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
    // Media container, list, header.
    var mediaCont = document.getElementById("media-container");
    var media = document.getElementById("media");
    var mediaHead = document.getElementById("media-header");

    // Play list container, list, and header.
    var playlistCont = document.getElementById("playlist-container");
    var playlists = document.getElementById("playlists");
    var playlistHead = document.getElementById("playlist-header");

    // Upper player container and lower "GOPF" label.
    var playerCont = document.getElementById("player-container");
    var footer = document.getElementById("footer");

    var width = window.innerWidth - playlistCont.offsetLeft -
                playlistCont.offsetWidth - DIVIDER_WIDTH;
    var height = footer.offsetTop - playerCont.offsetHeight -
                 playerCont.offsetTop - playlistCont.offsetLeft;

    mediaCont.style.maxWidth = width + "px";
    mediaCont.style.maxHeight = height + "px";
    media.style.maxHeight = height - mediaHead.offsetHeight + "px";

    playlistCont.style.maxHeight = height + "px";
    playlists.style.maxHeight = height - playlistHead.offsetHeight + "px";
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
          select(document.getElementById("media").childNodes[0]);
        }
      } else {
        select(document.getElementById("playlists").childNodes[0]);
      }
    }
  };
})();

window.addEventListener("load", input.init, true);
