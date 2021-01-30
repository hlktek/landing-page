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
  }, 3600000); // 3 hours = 10800000 miliseconds

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

  function formatNumber(number) {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND',
    }).format(number)
  }

  $(".nav-top-winner").on("click", function (e) {
    $("#top-winners-nav").find('.nav-link').removeClass('active')
    $(this).find('.nav-link').addClass('active')
    var cat = $(this).data("cat");
    var htmlStrTopWinner = "";
    fetch('/getTopWinner?category='+cat)
    .then(response => response.json())
    .then(res => {
      if(res && res.data && res.data.length>0){
        res.data.forEach(element => {
            htmlStrTopWinner += `<li class="list-group-item">
                        <div class="media">
                            <div class="media-body">
                                <div class="jp-winner-name">`+element.displayName+`</div>
                            </div>
                            <div class="jp-winner-total">
                            `+formatNumber(element.totalWin)+`
                            </div>
                            
                        </div>
                    </li>`
        });
        $("#top-winners-body").html(htmlStrTopWinner)
      }
    });
  });

});
