var loc = (function() {
  function deletePage() {
    var body = document.getElementsByTagName("body")[0];
    body.innerHTML = "<div style=\"text-align:center;margin-top:40px\"><h1> NO SOUP FOR YOU </h1></div>"; }

  function errorHandler() {
    window.console.log("Therer was an error.");
    deletePage();
  }

  function positionHandler(pos) {
    logPos(loc.pos = pos);
  }

  return {
    pos : {},

    init : function(e) {
      if (!navigator.geolocation) {
        deletePage();
      }

      navigator.geolocation.getCurrentPosition(positionHandler, errorHandler);
    }
  };
})();

window.addEventListener("load", loc.init, true);
