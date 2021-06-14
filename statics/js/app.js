$(document).ready(function () {
  // 初始化 materialize
  M.AutoInit();

  // 设置导航栏激活状态
  $("#nav_stock").click(function () {
    $(this).addClass("active");
    $("nav li").removeClass("active");
  });

  $("#nav_fund").click(function () {
    $(this).addClass("active");
    $("nav li").removeClass("active");
  });

  // 筛选表单中开关显示检测表单
  $("#selector_with_checker").click(function () {
    $("#checker_options").toggle();
  });

  // 表单提交按钮点击事件
  $("#selector_submit_btn").click(function () {
    $(this).addClass("disabled");
    $("#model_header").text($(this).text() + "中，请稍候...");
    $("#load_modal").modal()[0].M_Modal.options.dismissible = false;
    $("#load_modal").modal("open");
    $.ajax({
      url: "/selector",
      type: "post",
      data: $("#selector_form").serialize(),
      success: function (data) {
        if (data.Error != "") {
          $("#err_msg").text(data.Error);
          $("#error_modal").modal("open");
          $("#selector_submit_btn").removeClass("disabled");
          $("#load_modal").modal("close");
          return;
        }
        $.each(data.Stocks, function (i, stock) {
          $("#selector_result tbody").append(
            "<tr>" +
              "<td>" +
              stock.code.split(".")[0] +
              "</td>" +
              '<td><a href="#!">' +
              stock.name +
              "</a></td>" +
              "</tr>"
          );
        });
        $("title").text(data.PageTitle);
        $("#stock_forms").remove();
        $("#selector_result").removeClass("hide");
        $("#load_modal").modal("close");
      },
    });
  });

  $("#checker_submit_btn").click(function () {
    if ($("#checker_keyword").val() == "") {
      $("#err_msg").text("请填写股票代码或简称");
      $("#error_modal").modal("open");
      return;
    }
    $(this).addClass("disabled");
    $("#model_header").text($(this).text() + "中，请稍候...");
    $("#load_modal").modal()[0].M_Modal.options.dismissible = false;
    $("#load_modal").modal("open");
    $.ajax({
      url: "/checker",
      type: "post",
      data: $("#checker_form").serialize(),
      success: function (data) {
        if (data.Error != "") {
          $("#err_msg").text(data.Error);
          $("#error_modal").modal("open");
          $("#checker_submit_btn").removeClass("disabled");
          $("#load_modal").modal("close");
          return;
        }
        $.each(data.Results, function (i, result) {
          $("#checker_results").append(
            '<div class="divider"></div>' +
              '<div id="checker_result_' +
              i +
              '" class="row">' +
              '<div class="row"><h6>' +
              data.Names[i] +
              "</h6>" +
              '<table class="striped">' +
              '<thead><tr><th width="25%">指标</th><th width="65%">描述</th><th width="10%">结果</th></tr></thead>' +
              "<tbody></tbody>" +
              "</table>" +
              "</div>" +
              "</div>"
          );
          $.each(result, function (k, v) {
            okdesc = "❌异常";
            if (v.ok == "true") {
              okdesc = "✅正常";
            }
            $("#checker_result_" + i + " tbody").append(
              "<tr><td>" +
                k +
                "</td><td>" +
                v.desc +
                "</td><td>" +
                okdesc +
                "</td></tr>"
            );
          });
        });
        $("title").text(data.PageTitle);
        $("#stock_forms").remove();
        $("#checker_results").removeClass("hide");
        $("#load_modal").modal("close");
      },
    });
  });
});
