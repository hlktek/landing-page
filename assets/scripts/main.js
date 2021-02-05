$(document).ready(function () {
  var $newGameGrid = $("#new-game-grid");
  var $galleryGrid = $("#list-game");

  if ($.fn.isotope) {
    // $newGameGrid.imagesLoaded(function () {
    //   $newGameGrid.isotope({ itemSelector: ".item", layoutMode: "fitRows" });
    // });

    // $galleryGrid.imagesLoaded(function () {
    //   $galleryGrid.isotope({
    //     itemSelector: ".item",
    //     layoutMode: "fitColumns",
    //     percentPosition: true,
    //     // masonry: {
    //     //   columnWidth: ".grid-sizer",
    //     // },
    //   });
    // });

    // var $grid = $galleryGrid.isotope({ itemSelector: ".item", layoutMode: "masonry" });
    var $gridSelectors = $(".game-filter").find("a");

    $gridSelectors.on("click", function (e) {
      var selector = $(this).attr("data-filter");
      $galleryGrid.find(".item").not(selector).hide();
      $galleryGrid.find(selector).fadeIn();
      // console.log(selector);
      // $galleryGrid.isotope({
      //   filter: selector,
      //   layoutMode: "fitRows",
      // });

      $gridSelectors.removeClass("actived");
      $(this).addClass("actived");

      e.preventDefault();
    });

    var qsRegex;

    $("#searchTerm").keyup(
      debounce(function () {
        qsRegex = new RegExp($("#searchTerm").val(), "gi");

        $galleryGrid.isotope({
          itemSelector: ".item",
          filter: function () {
            return qsRegex ? $(this).text().match(qsRegex) : true;
          },
        });
        e.preventDefault();
      }, 200)
    );
  }

  $("#slide-banner").carousel();

  $("#logout").on("click", function (e) {
    clearCookies("mysession");
  });

  setInterval(function () {
    clearCookies("mysession");
  }, 3600000); // 3 hours = 10800000 miliseconds

  // getJackpotHistory
  getJackpotHistory();

  setInterval(function () {
    getJackpotHistory();
  }, 300000); // 5 phut

  // get top winner
  getTopWinner("all");
  setInterval(function () {
    getTopWinner("all");
  }, 300000); // 5 phut

  $(".nav-top-winner").on("click", function (e) {
    $("#top-winners-nav").find(".nav-link").removeClass("active");
    $(this).find(".nav-link").addClass("active");
    var cat = $(this).data("cat");
    getTopWinner(cat);
  });

  function clearCookies(cookieName) {
    var cookies = document.cookie.split(";");

    for (var i = 0; i < cookies.length; i++) {
      var cookie = cookies[i];
      var eqPos = cookie.indexOf("=");
      var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;

      if (name.trim() == cookieName) {
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
        window.location.reload();
      }
    }
  }

  function formatNumber(number) {
    return new Intl.NumberFormat("vi-VN", {
      style: "currency",
      currency: "VND",
    }).format(number);
  }

  function getTopWinner(cat) {
    $("#top-winners-body").html('<div class="spinner-border spin-loading" role="status"></div>');
    fetch("/getTopWinner?category=" + cat)
      .then((response) => response.json())
      .then((res) => {
        var htmlStrTopWinner = "";
        if (res && res.data && res.data.length > 0) {
          res.data.forEach((element) => {
            htmlStrTopWinner +=
              `<li class="list-group-item">
                                    <div class="media">
                                        <div class="media-body">
                                            <div class="jp-winner-name">` +
              element.displayName +
              `</div>
                                        </div>
                                        <div class="jp-winner-total">
                                        ` +
              formatNumber(element.totalWinLoss) +
              `
                                        </div>

                                    </div>
                                </li>`;
          });
          $("#top-winners-body").html(htmlStrTopWinner);
        } else {
          $("#top-winners-body").html("<div class='text-center mt-3'>Không có data</div>");
        }
      })
      .catch(function () {
        $("#top-winners-body").html("<div class='text-center mt-3'>Có lỗi xin vui lòng thử lại</div>");
      });
  }

  function getJackpotHistory() {
    fetch("/getJackpotHistory")
      .then((response) => response.json())
      .then((res) => {
        var htmlStrTopWinner = "";
        if (res && res.data && res.data.length > 0) {
          res.data.forEach((element) => {
            htmlStrTopWinner +=
              `<li class="list-group-item">
                                    <div class="media">
                                        <div class="media-body">
                                            <div class="jp-winner-game">` +
              element.displayName +
              `</div>
                                            <div class="jp-game-name">` +
              element.serviceId +
              `</div>
                                        </div>
                                        <div class="jp-total">
                                          ` +
              formatNumber(element.jackpotAmount) +
              `
                                        </div>
                                        
                                    </div>
                                </li>`;
          });
          $("#jp-body").html(htmlStrTopWinner);
        } else {
          $("#jp-body").html("<div class='text-center mt-3'>Không có data</div>");
        }
      })
      .catch(function () {
        $("#jp-body").html("<div class='text-center mt-3'>Có lỗi xin vui lòng thử lại</div>");
      });
  }

  // debounce so filtering doesn't happen every millisecond
  function debounce(fn, threshold) {
    var timeout;
    threshold = threshold || 100;
    return function debounced() {
      clearTimeout(timeout);
      var args = arguments;
      var _this = this;
      function delayed() {
        fn.apply(_this, args);
      }
      timeout = setTimeout(delayed, threshold);
    };
  }
});
