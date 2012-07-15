// Copyright 2012 Andrew "Jamoozy" Correa
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License
// for more details.
//
// You should have received a copy of the GNU General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.

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
    var player = document.getElementById("player");
    switch(e.keyCode) {
      case 33:  // PageUp
        prev(10);
        e.stopPropagation();
        break;
      case 75:  // k
        prev(1);
        e.stopPropagation();
        break;
      case 34:  // PageDown
        next(10);
        e.stopPropagation();
        break;
      case 74:  // j
        next(1);
        e.stopPropagation();
        break;
      case 192: // back-tick ("`")
      case 72:  // h
        input.swap();
        e.stopPropagation();
        break;

      case 38:  // up
        player.currentTime += e.ctrlKey ? 600 : 300;
        break;
      case 40:  // down
        player.currentTime -= e.ctrlKey ? 600 : 300;
        break;
      case 37:  // left
        player.currentTime -= e.ctrlKey ? 60 : 10;
        break;
      case 39:  // right
        player.currentTime += e.ctrlKey ? 60 : 10;
        break;

      case 80:  // P
        media.prev();
        break;
      case 78:  // N
        media.next();
        break;

      case 83:  // S
        document.getElementById("shuf").click();
        break;
      case 76:  // L
        document.getElementById("loop").click();
        break;

      case 219:  // [
        player.playbackRate -= 0.5;
        break;
      case 221:  // ]
        player.playbackRate += 0.5;
        break;
      case 8:  // backspace
        player.playbackRate = 1.0;
        break;

      case 48:  // 0
        player.removeAttribute("width");
        adjustSize();
        break;
      case 49:  // 1
        player.width = 200;
        adjustSize();
        break;
      case 50:  // 2
        player.width = 400;
        adjustSize();
        break;
      case 51:  // 3
        player.width = 600;
        adjustSize();
        break;
      case 52:  // 4
        player.width = 800;
        adjustSize();
        break;
      case 53:  // 5
        player.width = 1000;
        adjustSize();
        break;

      case 77:  // M
        player.muted = !player.muted;
        break;
      case 32:  // space
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

      // TODO key for displaying help
      case 191:  // ?
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
    var controls = document.getElementById("controls");
    var anElem = document.getElementsByClassName("unselected")[0];

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

    // Set max dims of the video element.
    var margin = 10; // px
    player.style.maxWidth = playerCont.offsetWidth - 2 * margin + "px"
    var topSp = player.offsetTop;
    var bottomSp = mediaHead.offsetHeight + controls.offsetHeight + 3 * anElem.offsetHeight;
    player.style.maxHeight = footer.offsetTop - topSp - bottomSp + "px";
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

      window.addEventListener("resize", function(e) { adjustSize(); }, true);
      window.addEventListener("keydown", onkey, true);
      document.getElementById("player").addEventListener(
          "canplay", function(e) { adjustSize(); }, true);
      adjustSize();
    },

    // Swaps the list the user is moving through.
    swap : function() {
      viewMedia = !viewMedia;
      if (viewMedia) {
        if (document.getElementsByClassName("dummy").length > 0) {
          viewMedia = false;
        } else {
          var media = document.getElementById("media");
          var playing = media.getElementsByClassName("playing");
          if (playing.length > 0) {
            select(playing[0]);
          } else {
            select(media.childNodes[0]);
          }
        }
      } else {
        var playlists = document.getElementById("playlists");
        var selected = playlists.getElementsByClassName("selected");
        if (selected.length > 0) {
          select(selected[0]);
        } else {
          select(playlists.childNodes[0]);
        }
      }
    }
  };
})();

window.addEventListener("load", input.init, true);
