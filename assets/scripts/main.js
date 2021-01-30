$(document).ready(function () {
  var $galleryGrid = $(".game-grid");

  if ($.fn.isotope) {
    $galleryGrid.imagesLoaded(function () {
      $galleryGrid.isotope({ itemSelector: ".item", layoutMode: "masonry" });
    });

    // var $grid = $galleryGrid.isotope({ itemSelector: ".item", layoutMode: "masonry" });
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
    clearCookies('mysession');
  });
  
  setInterval(function(){
    console.log("Hello"); 
    clearCookies('mysession');
  }, 10800000); // 3 hours = 10800000 miliseconds

  function clearCookies(cookieName) {
    var cookies = document.cookie.split(";");

    for (var i = 0; i < cookies.length; i++) {
      var cookie = cookies[i];
      var eqPos = cookie.indexOf("=");
      var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
      
      if(name.trim() == cookieName){
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
        window.location.reload();
      }
    }
  }

});
