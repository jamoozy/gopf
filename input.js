// Copyright 2013 Andrew "Jamoozy" Correa S.
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
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

  var helpOrder = [
    'Navigation',
    33, 34, 75, 74, 72, 13,
    'Scanning',
    38, 40, 37, 39,
    'Player Controls',
    80, 78, 83, 76, 77, 32,
    'Speed Controls',
    219, 221, 8,
    'Size Controls',
    48, 49, 50, 51, 52, 53,
    'Other',
    191
  ];


  function bindings() {
    var player = document.getElementById("player");
    return {
      33: {
        key: 'Pg&uarr;',
        use: 'Move up 10.',
        func: function(e) {
          prev(10);
          e.stopPropagation();
        }
      },
      34: {
        key: 'Pg&darr;',
        use: 'Move down 10.',
        func: function(e) {
          next(10);
          e.stopPropagation();
        }
      },
      75: {
        key: 'K',
        use: 'Move up 1.',
        func: function(e) {
          prev(1);
          e.stopPropagation();
        }
      },
      74: {
        key: 'J',
        use: 'Move down 1.',
        func: function(e) {
          next(1);
          e.stopPropagation();
        }
      },
      72: {
        key: 'H',
        use: 'Switch between media and playlist lists.',
        func: function(e) {
          input.swap();
          e.stopPropagation();
        }
      },
      13: {
        key: 'Enter',
        use: 'Select highlighted playlist/media.',
        func: function(e) {
          var sel = document.getElementById("selected");
          if (sel !== null && sel !== undefined) {
            if (viewMedia) {
              media.onclick(sel);
            } else {
              playlist.onclick(sel, true);
            }
          }
        }
      },

      38: {
        key: '&uarr; / Ctrl &uarr;',
        use: 'Go forward 5 / 10 minutes',
        func: function(e) {
          player.currentTime += e.ctrlKey ? 600 : 300;
        }
      },
      40: {
        key: '&darr; / Ctrl &darr;',
        use: 'Go back 5 / 10 minutes',
        func: function(e) {
          player.currentTime -= e.ctrlKey ? 600 : 300;
        }
      },
      37: {
        key: '&larr; / Ctrl &larr;',
        use: 'Go forward 10 / 60 seconds',
        func: function(e) {
          player.currentTime -= e.ctrlKey ? 60 : 10;
        }
      },
      39: {
        key: '&rarr; / Ctrl &rarr;',
        use: 'Go back 10 / 60 seconds',
        func: function(e) {
          player.currentTime += e.ctrlKey ? 60 : 10;
        }
      },

      80: {
        key: 'P',
        use: 'Previous track',
        func : function(e) {
          media.prev();
        }
      },
      78: {
        key: 'N',
        use: 'Next track',
        func : function(e) {
          media.next();
        }
      },
      83: {
        key: 'S',
        use: 'Toggle Shuffle',
        func: function(e) {
          document.getElementById("shuf").click();
        }
      },
      76: {
        key: 'L',
        use: 'Toggle Loop',
        func: function(e) {
          document.getElementById("loop").click();
        }
      },
      77: {
        key: 'M',
        use: 'Mute/unmute player',
        func: function(e) {
          player.muted = !player.muted;
        }
      },
      32: {
        key: 'Spbar',
        use: 'Pause / unpause',
        func: function(e) {
          if (player.paused) {
            player.play();
          } else {
            player.pause();
          }
        }
      },

      219: {
        key: '[',
        use: 'Subtract 0.5 from playback rate',
        func: function(e) {
          player.playbackRate -= 0.5;
        }
      },
      221: {
        key: ']',
        use: 'Add 0.5 to playback rate',
        func: function(e) {
          player.playbackRate += 0.5;
        }
      },
      8: {
        key: 'Bksp',
        use: 'Return playback to normal speed',
        func: function(e) {
          player.playbackRate = 1.0;
        }
      },

      48: {
        key: '0',
        use: 'Set video width to max',
        func: function(e) {
          player.removeAttribute("width");
          adjustSize();
        }
      },
      49: {
        key: '1',
        use: 'Set video width to 200px',
        func: function(e) {
          player.width = 200;
          adjustSize();
        }
      },
      50: {
        key: '2',
        use: 'Set video width to 400px',
        func: function(e) {
          player.width = 400;
          adjustSize();
        }
      },
      51: {
        key: '3',
        use: 'Set video width to 600px',
        func: function(e) {
          player.width = 600;
          adjustSize();
        }
      },
      52: {
        key: '4',
        use: 'Set video width to 800px',
        func: function(e) {
          player.width = 800;
          adjustSize();
        }
      },
      53: {
        key: '5',
        use: 'Set video width to 1000px',
        func: function(e) {
          player.width = 1000;
          adjustSize();
        }
      },

      191: {
        key: '?',
        use: 'Open / close this help dialog.',
        func: function(e) {
          toggleHelpDialog();
        }
      }
    };
  }

  function toggleHelpDialog() {
    var hd = document.getElementById("help-dialog");
    if (hd.style.visibility == "hidden") {
      hd.style.visibility = "visible";
    } else {
      hd.style.visibility = "hidden";
    }
  }

  function initHelpDialog() {
    var b = bindings();
    var html = "<ul>";
    for (var i in helpOrder) {
      var elem = helpOrder[i];
      if (typeof elem === "string") {
        html += '<li class="help-header">' + elem + "</li>";
      } else {
        window.console.log("Checking out element: " + elem);
        html += '<li><span class="key">' + b[elem].key + '</span>: ' +
          '<span class="use">' + b[elem].use + "</span></li>";
      }
    }
    html += "</ul>";
    document.getElementById('help-dialog').innerHTML = html;
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
    var b = bindings()
    if (e.keyCode in b) {
      b[e.keyCode].func(e);
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
    var player = document.getElementById("player");
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

      document.getElementById("help-dialog").style.visibility = "hidden";
      window.addEventListener("resize", function(e) { adjustSize(); }, true);
      window.addEventListener("keydown", onkey, true);
      document.getElementById("player").addEventListener(
          "canplay", function(e) { adjustSize(); }, true);
      initHelpDialog();
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
