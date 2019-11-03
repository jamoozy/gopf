// Copyright (C) 2011-2019 Andrew "Jamoozy" C. Sabisch
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

var input = (function() {
  // Extra amount to shrink media container by, so that media/playlists
  // don't overlap one onther.
  var DIVIDER_WIDTH = 40;

  // Whether the selector is in the media list.  False means it's in the
  // playlists list.
  var viewMedia = false;

  // Fullscreened (in a window) state.
  var fullscreen = false;

  function notify(str) {
    $("#notification").html(str);
  }

  var helpOrder = [
    'Navigation',
    33, 34, 75, 74, 72, 13,
    'Scanning',
    38, 40, 37, 39,
    'Player Controls',
    70, 80, 78, 83, 76, 77, 32,
    'Speed Controls',
    219, 221, 8,
    'Size Controls',
    48, 49, 50, 51, 52, 53,
    'Other',
    191
  ];


  function bindings() {
    var player = $("#player");
    return {
      33: {
        key: 'Pg&uarr;',
        use: 'Move up 10.',
        func: function(e) {
          move(-10);
          return true;
        }
      },
      34: {
        key: 'Pg&darr;',
        use: 'Move down 10.',
        func: function(e) {
          move(10);
          return true;
        }
      },
      75: {
        key: 'K',
        use: 'Move up 1.',
        func: function(e) {
          move(-1);
          return true;
        }
      },
      74: {
        key: 'J',
        use: 'Move down 1.',
        func: function(e) {
          move(1);
          return true;
        }
      },
      72: {
        key: 'H',
        use: 'Switch between media and playlist lists.',
        func: function(e) {
          input.swap();
          return true;
        }
      },
      13: {
        key: 'Enter',
        use: 'Select highlighted playlist/media.',
        func: function(e) {
          var sel = $("#selected");
          if (sel !== null && sel !== undefined) {
            if (viewMedia) {
              $(media).trigger('click', sel);
            } else {
              $(playlist).trigger('click', sel, true);
            }
            return true;
          }
        }
      },

      38: {
        key: '&uarr; / Ctrl &uarr;',
        use: 'Go forward 5 / 10 minutes',
        func: function(e) {
          player[0].currentTime += e.ctrlKey ? 600 : 300;
          return true;
        }
      },
      40: {
        key: '&darr; / Ctrl &darr;',
        use: 'Go back 5 / 10 minutes',
        func: function(e) {
          player[0].currentTime -= e.ctrlKey ? 600 : 300;
          return true;
        }
      },
      37: {
        key: '&larr; / Ctrl &larr;',
        use: 'Go forward 10 / 60 seconds',
        func: function(e) {
          player[0].currentTime -= e.ctrlKey ? 60 : 10;
          return true;
        }
      },
      39: {
        key: '&rarr; / Ctrl &rarr;',
        use: 'Go back 10 / 60 seconds',
        func: function(e) {
          player[0].currentTime += e.ctrlKey ? 60 : 10;
          return true;
        }
      },

      70: {
        key: 'F',
        use: 'Toogle (windowed) fullscreen mode',
        fun: function(e) {
          toggleFullscreen();
          return true;
        }
      },
      80: {
        key: 'P',
        use: 'Previous track',
        func : function(e) {
          media.prev();
          return true;
        }
      },
      78: {
        key: 'N',
        use: 'Next track',
        func : function(e) {
          media.next();
          return true;
        }
      },
      83: {
        key: 'S',
        use: 'Toggle Shuffle',
        func: function(e) {
          $("#shuf").trigger('click');
          return true;
        }
      },
      76: {
        key: 'L',
        use: 'Toggle Loop',
        func: function(e) {
          if (!e.ctrlKey && !e.metaKey) {
            $("#loop").trigger('click');
            return true;
          }
        }
      },
      77: {
        key: 'M',
        use: 'Mute/unmute player',
        func: function(e) {
          player[0].muted = !player[0].muted;
          return true;
        }
      },
      32: {
        key: 'Spbar',
        use: 'Pause / unpause',
        func: function(e) {
          if (player[0].paused) {
            player[0].play();
          } else {
            player[0].pause();
          }
          return true;
        }
      },

      219: {
        key: '[',
        use: 'Subtract 0.5 from playback rate',
        func: function(e) {
          player[0].playbackRate -= 0.5;
          return true;
        }
      },
      221: {
        key: ']',
        use: 'Add 0.5 to playback rate',
        func: function(e) {
          player[0].playbackRate += 0.5;
          return true;
        }
      },
      8: {
        key: 'Bksp',
        use: 'Return playback to normal speed',
        func: function(e) {
          player[0].playbackRate = 1.0;
          return true;
        }
      },

      48: {
        key: '0',
        use: 'Set video width to max',
        func: function(e) {
          if (e.ctrlKey || e.metaKey) {
            return false;
          }
          player.removeAttr("width");
          adjustSize();
          return true;
        }
      },
      49: {
        key: '1',
        use: 'Set video width to 200px',
        func: function(e) {
          if (e.ctrlKey || e.metaKey) {
            return false;
          }
          player.attr('width', 200);
          adjustSize();
          return true;
        }
      },
      50: {
        key: '2',
        use: 'Set video width to 400px',
        func: function(e) {
          if (e.ctrlKey || e.metaKey) {
            return false;
          }
          player.attr('width', 400);
          adjustSize();
          return true;
        }
      },
      51: {
        key: '3',
        use: 'Set video width to 600px',
        func: function(e) {
          if (e.ctrlKey || e.metaKey) {
            return false;
          }
          player.attr('width', 600);
          adjustSize();
          return true;
        }
      },
      52: {
        key: '4',
        use: 'Set video width to 800px',
        func: function(e) {
          if (e.ctrlKey || e.metaKey) {
            return false;
          }
          player.attr('width', 800);
          adjustSize();
          return true;
        }
      },
      53: {
        key: '5',
        use: 'Set video width to 1000px',
        func: function(e) {
          if (e.ctrlKey || e.metaKey) {
            return false;
          }
          player.attr('width', 1000);
          adjustSize();
          return true;
        }
      },

      191: {
        key: '?',
        use: 'Open / close this help dialog.',
        func: function(e) {
          toggleHelpDialog();
          return true;
        }
      }
    };
  }

  function toggleFullscreen() {
    var player = $("#player");
    if (player[0].tagName === 'AUDIO') {
      return;
    }
    if (fullscreen = !fullscreen) {
    } else {
    }
  }

  function toggleHelpDialog() {
    var hd = $("#help-dialog");
    if (hd.css("display") === "none") {
      hd.show();
    } else {
      hd.hide();
    }
  }

  function initHelpDialog() {
    var b = bindings();
    var html = "<ul>";
    $(helpOrder).each(function(i, elem) {
      if (typeof elem === "string") {
        html += '<li class="help-header">' + elem + "</li>";
      } else {
        html += '<li><span class="key">' + b[elem].key + '</span>: ' +
          '<span class="use">' + b[elem].use + "</span></li>";
      }
    });
    html += "</ul>";
    $("#help-dialog").html(html);
    //$("#help-dialog").hide();
  }

  // Selects the given element.
  function select(idx) {
    var sel = $("#selected").removeAttr("id");
    $(getElems().get(idx)).attr("id", "selected");
  }

  // Gets the elements from the currently-selected list.
  function getElems() {
    return getList().children();
  }

  // Gets the currently-selected list.
  function getList() {
    return viewMedia ? $("#media") : $("#playlists");
  }

  // Ensures the selected element is visible, either by centering the
  // list's view on the element, or scrolling the list all the way up or
  // down so the element is in view.
  function ensureSelectedVisible(i) {
    var list = getList();
    // "-1" due to "border-bottom: -1px"
    list.scrollTop(i * ($("#selected").innerHeight() + 1) - list.height() / 2);
  }

  function getSelectedIndex() {
    var elems = getElems();
    for (var i = 0; i < elems.size(); i++) {
      if ($(elems.get(i)).attr('id') === 'selected') {
        return i;
      }
    }
  }

  // Selects the next element in the list.
  function move(inc) {
    if (inc === 0) {
      return;
    }

    var elems = getElems();
    var i = getSelectedIndex();
    var val = i + inc;
    if (val >= elems.size()) {
      val = elems.size() - 1;
    } else if (val < 0) {
      val = 0;
    }

    select(val);
    ensureSelectedVisible(val);
  }

  // Key handler.
  function onkey(e) {
    var b = bindings();
    if (e.keyCode in b && b[e.keyCode].func(e)) {
      e.preventDefault();
    }
  }

  // Window resize handler.
  function adjustSize() {
    // Media container, list, header.
    var mediaCont = $("#media-container");
    var media = $("#media");
    var mediaHead = $("#media-header");

    // Play list container, list, and header.
    var playlistCont = $("#playlist-container");
    var playlists = $("#playlists");
    var playlistHead = $("#playlist-header");
    var controls = $("#controls");
    var anElem = $(".unselected").first();

    // Upper player container and lower "GOPF" label.
    var playerCont = $("#player-container");
    var player = $("#player");
    var footer = $("#footer");

    var width = $(window).innerWidth() - playlistCont.offset().left -
                playlistCont.outerWidth(true) - DIVIDER_WIDTH;
    var height = footer.offset().top - playerCont.outerHeight(true) -
                 playerCont.offset().top - playlistCont.offset().left;

    mediaCont.css("max-width", width + "px");
    mediaCont.css("min-width", width + "px");
    mediaCont.css("max-height", height + "px");
    mediaCont.css("min-height", height + "px");
    media.css("max-height", height - mediaHead.outerHeight(true) + "px");
    media.css("min-height", height - mediaHead.outerHeight(true) + "px");

    playlistCont.css("max-height", height + "px");
    playlistCont.css("min-height", height + "px");
    playlists.css("max-height", height - playlistHead.outerHeight(true) + "px");
    playlists.css("min-height", height - playlistHead.outerHeight(true) + "px");

    // Set max dims of the video element.
    var margin = 10; // px
    player.css("width", playerCont.outerWidth(true) - 2 * margin + "px");
  }

  return {
    // Initializes the input module by registering event listeners.
    init : function() {
      var sel = $("#selected")[0];
      if (!sel) {
        // Default selected element is first playlist.
        $("#playlists").children().first().attr("id", "selected");
      }

      $(window).resize(function(e) { adjustSize(); });
      $(window).keydown(onkey);
      $("#help-dialog").hide();
      $("#player").on("canplay", function(e) { adjustSize(); });
      initHelpDialog();
      adjustSize();
    },

    // Swaps the list the user is moving through.
    swap : function() {
      viewMedia = !viewMedia;
      if (viewMedia) {
        if ($(".dummy").length > 0) {
          viewMedia = false;
        } else {
          select(0);
        }
      } else {
        select(0);
      }
    }
  };
})();

$(window).load(input.init);
