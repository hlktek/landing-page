$(document).ready(function () {
  var $galleryGrid = $(".game-grid");
  if ($.fn.isotope) {
    $galleryGrid.isotope({ itemSelector: ".item", layoutMode: "masonry" });
    var $gridSelectors = $(".game-filter").find("a");

    $gridSelectors.on("click", function (e) {
      var selector = $(this).attr("data-filter");

      $galleryGrid.isotope({
        filter: selector,
      });

      $gridSelectors.removeClass("actived");
      $(this).addClass("actived");

      e.preventDefault();
    });
  }

  $("#slide-banner").carousel();

  $("#logout").on("click", function (e) {
    var cookies = document.cookie.split(";");

    for (var i = 0; i < cookies.length; i++) {
      var cookie = cookies[i];
      var eqPos = cookie.indexOf("=");
      var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
      document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
    }
    window.location.reload();
  });
  
});
