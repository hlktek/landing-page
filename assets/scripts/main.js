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

  $("#btn-naptien").on("click", function (e) {
    var _self = this;
    $(_self).html('<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>');
    $(_self).attr("disabled", "disabled");
    fetch("/addWallet")
      .then((response) => response.json())
      .then((res) => {
        $("#modalNaptien").modal("hide");
        $(_self).html("Đồng ý");
        $(_self).removeAttr("disabled", "disabled");
        if (res.data && res.data.code == 200) {
          toastr.success("Nạp tiền thành công");
        } else {
          toastr.error(res.data && res.data.message);
        }
      })
      .catch(function (err) {
        $(_self).html("Đồng ý");
        $(_self).removeAttr("disabled", "disabled");
        $("#modalNaptien").modal("hide");
        toastr.error("Nạp tiền thất bại, xin vui lòng thử lại.");
      });
  });

  $("#btn-show-confirm").on("click", function (e) {
    var feedBack = $("#feedBack").val();
    if (feedBack.trim() === "") {
      toastr.error("Vui lòng nhập đủ thông tin");
    } else {
      $("#modalFeedback").modal("show");
    }
  });

  $("#btn-send-feedback").on("click", function (e) {
    $("#frm-feedback").submit();
  });

  $("#frm-feedback").submit(function (event) {
    var userId = $("#email").val();
    var serviceId = $("#serviceId").val();
    var feedBack = $("#feedBack").val();

    if (feedBack.trim() === "") {
      toastr.error("Vui lòng nhập đủ thông tin");
    } else {
      $("#btn-send-feedback").html(
        '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>'
      );
      $("#btn-send-feedback").attr("disabled", "disabled");
      var data = {
        userId: userId,
        serviceId: serviceId,
        feedBack: feedBack,
      };
      fetch("/insertFeedbackEs", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      })
        .then((response) => response.json())
        .then((res) => {
          $("#modalFeedback").modal("hide");
          $("#btn-send-feedback").html("Đồng ý");
          $("#btn-send-feedback").removeAttr("disabled", "disabled");
          if (res.code === 200) {
            toastr.success(res.message);
            $("#feedBack").val("");
          } else {
            toastr.error("Gửi feedback thất bại, xin vui lòng thử lại.");
          }
        })
        .catch((error) => {
          $("#modalFeedback").modal("hide");
          $("#btn-send-feedback").html("Đồng ý");
          $("#btn-send-feedback").removeAttr("disabled", "disabled");
          toastr.error("Gửi feedback thất bại, xin vui lòng thử lại.");
        });
    }
    event.preventDefault();
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
