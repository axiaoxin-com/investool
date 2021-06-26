$(document).ready(function () {
  // 初始化 materialize
  M.AutoInit();

  // 导航栏激活
  var currentNav = $(location).attr("pathname");
  if (currentNav == "/fund") {
    $("#nav-fund").addClass("active");
    $("#nav-fund").siblings().removeClass("active");
  } else if (currentNav == "/about") {
    $("#nav-about").addClass("active");
    $("#nav-about").siblings().removeClass("active");
  } else if (currentNav == "/comment") {
    $("#nav-comment").addClass("active");
    $("#nav-comment").siblings().removeClass("active");
  } else {
    $("#nav-stock").addClass("active");
    $("#nav-stock").siblings().removeClass("active");
  }

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
        if (data.Stocks.length == 0) {
          $("#dropdown-structure").addClass("hide");
          $("#selector_result #result_table").html(
            '<div class="row"><p class="center flow-text">无法找到符合条件的股票</p></div>'
          );
        } else {
          $.each(data.Stocks, function (i, stock) {
            var cm = stock.code.split(".");
            $("#selector_result tbody").append(
              "<tr>" +
                "<td>" +
                cm[0] +
                "</td>" +
                '<td><a target="_blank" href="http://quote.eastmoney.com/' +
                cm[1] +
                cm[0] +
                '.html">' +
                stock.name +
                "</a></td>" +
                '<td class="hide t_industry">' +
                stock.industry +
                "</td>" +
                '<td class="hide t_keywords">' +
                stock.keywords +
                "</td>" +
                '<td class="hide t_company_profile">' +
                stock.company_profile +
                "</td>" +
                '<td class="hide t_main_forms">' +
                stock.main_forms +
                "</td>" +
                '<td class="hide t_byys_ration">' +
                stock.byys_ration +
                "</td>" +
                '<td class="hide t_report_date_name">' +
                stock.report_date_name +
                "</td>" +
                '<td class="hide t_report_opinion">' +
                stock.report_opinion +
                "</td>" +
                '<td class="hide t_jzpg">' +
                stock.jzpg +
                "</td>" +
                '<td class="hide t_latest_roe">' +
                stock.latest_roe +
                "</td>" +
                '<td class="hide t_roe_tbzz">' +
                stock.roe_tbzz +
                "</td>" +
                '<td class="hide t_roe_5y">' +
                stock.roe_5y +
                "</td>" +
                '<td class="hide t_latest_eps">' +
                stock.latest_eps +
                "</td>" +
                '<td class="hide t_eps_tbzz">' +
                stock.eps_tbzz +
                "</td>" +
                '<td class="hide t_eps_5y">' +
                stock.eps_5y +
                "</td>" +
                '<td class="hide t_total_income">' +
                stock.total_income +
                "</td>" +
                '<td class="hide t_total_income_tbzz">' +
                stock.total_income_tbzz +
                "</td>" +
                '<td class="hide t_total_income_5y">' +
                stock.total_income_5y +
                "</td>" +
                '<td class="hide t_net_profit">' +
                stock.net_profit +
                "</td>" +
                '<td class="hide t_net_profit_tbzz">' +
                stock.net_profit_tbzz +
                "</td>" +
                '<td class="hide t_net_profit_5y">' +
                stock.net_profit_5y +
                "</td>" +
                '<td class="hide t_zxgxl">' +
                stock.zxgxl +
                "</td>" +
                '<td class="hide t_fina_report_date">' +
                stock.fina_report_date +
                "</td>" +
                '<td class="hide t_fina_appoint_publish_date">' +
                stock.fina_appoint_publish_date +
                "</td>" +
                '<td class="hide t_fina_actual_publish_date">' +
                stock.fina_actual_publish_date +
                "</td>" +
                '<td class="hide t_total_market_cap">' +
                stock.total_market_cap +
                "</td>" +
                '<td class="hide t_price">' +
                stock.price +
                "</td>" +
                '<td class="hide t_right_price">' +
                stock.right_price +
                "</td>" +
                '<td class="hide t_price_space">' +
                stock.price_space +
                "</td>" +
                '<td class="hide t_hv">' +
                stock.hv +
                "</td>" +
                '<td class="hide t_zxfzl">' +
                stock.zxfzl +
                "</td>" +
                '<td class="hide t_fzldb">' +
                stock.fzldb +
                "</td>" +
                '<td class="hide t_netprofit_growthrate_3_y">' +
                stock.netprofit_growthrate_3_y +
                "</td>" +
                '<td class="hide t_income_growthrate_3_y">' +
                stock.income_growthrate_3_y +
                "</td>" +
                '<td class="hide t_listing_yield_year">' +
                stock.listing_yield_year +
                "</td>" +
                '<td class="hide t_listing_volatility_year">' +
                stock.listing_volatility_year +
                "</td>" +
                '<td class="hide t_pe">' +
                stock.pe +
                "</td>" +
                '<td class="hide t_peg">' +
                stock.peg +
                "</td>" +
                '<td class="hide t_org_rating">' +
                stock.org_rating +
                "</td>" +
                '<td class="hide t_profit_predict">' +
                stock.profit_predict +
                "</td>" +
                '<td class="hide t_valuation_syl">' +
                stock.valuation_syl +
                "</td>" +
                '<td class="hide t_valuation_sjl">' +
                stock.valuation_sjl +
                "</td>" +
                '<td class="hide t_valuation_sxol">' +
                stock.valuation_sxol +
                "</td>" +
                '<td class="hide t_valuation_sxnl">' +
                stock.valuation_sxnl +
                "</td>" +
                '<td class="hide t_hyjzsp">' +
                stock.hyjzsp +
                "</td>" +
                '<td class="hide t_ztzd">' +
                stock.ztzd +
                "</td>" +
                '<td class="hide t_mll_5y">' +
                stock.mll_5y +
                "</td>" +
                '<td class="hide t_jll_5y">' +
                stock.jll_5y +
                "</td>" +
                '<td class="hide t_listing_date">' +
                stock.listing_date +
                "</td>" +
                '<td class="hide t_netcash_operate">' +
                stock.netcash_operate +
                "</td>" +
                '<td class="hide t_netcash_invest">' +
                stock.netcash_invest +
                "</td>" +
                '<td class="hide t_netcash_finance">' +
                stock.netcash_finance +
                "</td>" +
                '<td class="hide t_netcash_free">' +
                stock.netcash_free +
                "</td>" +
                "</tr>"
            );
          });
        }
        $("title").text(data.PageTitle);
        $("#stock_forms").remove();
        $("#selector_result").removeClass("hide");
        $("html, body").animate({ scrollTop: 0 }, 0);
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
        if (data.Results.length == 0) {
          $("#checker_results h4").text("暂不支持对该股进行检测");
        } else {
          $.each(data.Results, function (i, result) {
            var cm = data.StockNames[i].split("-")[1].split(".");
            $("#checker_results").append(
              '<div class="divider"></div>' +
                '<div id="checker_result_' +
                i +
                '" class="row">' +
                '<div class="section"><h6><a class="right" target="_blank" href="http://quote.eastmoney.com/' +
                cm[1] +
                cm[0] +
                '.html">' +
                data.StockNames[i] +
                "</a><span>当前检测财报数据来源:" +
                data.FinaReportNames[i] +
                "</br></br>最新财报预约发布日期:" +
                data.FinaAppointPublishDates[i] +
                '</span></h6><div class="divider"></div>' +
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
        }
        $("title").text(data.PageTitle);
        $("#stock_forms").remove();
        $("#checker_results").removeClass("hide");
        $("html, body").animate({ scrollTop: 0 }, 0);
        $("#load_modal").modal("close");
      },
    });
  });

  // 返回顶部按钮
  $("#to-top").click(function () {
    $("html, body").animate({ scrollTop: 0 }, 500);
  });
  // 按钮通过点击展示
  $(".fixed-action-btn").floatingActionButton({
    hoverEnabled: false,
  });

  // 导出结果csv文件
  $(".export-result-btn").click(function (e) {
    tableExport("selector_result_table", "x-stock-selected", "csv");
  });

  // 展示字段设置
  var checkboxLimit = 10;
  var checkboxCountCheck = function () {
    var checkedCount = $("#fields-dropdown input[type=checkbox]:checked")
      .length;
    if (checkedCount > checkboxLimit && checkedCount % 5 == 1) {
      M.toast({
        html: "展示信息过多，导出CSV详情文件即可在本地查看完整信息哦~",
        classes: "rounded",
      });
    }
  };
  $("#fields-dropdown-btn").dropdown({
    constrainWidth: true,
    closeOnClick: false,
  });
  $("#field_industry").change(function () {
    checkboxCountCheck();
    $(".t_industry").toggleClass("hide");
  });
  $("#field_keywords").change(function () {
    checkboxCountCheck();
    $(".t_keywords").toggleClass("hide");
  });
  $("#field_company_profile").change(function () {
    checkboxCountCheck();
    $(".t_company_profile").toggleClass("hide");
  });
  $("#field_main_forms").change(function () {
    checkboxCountCheck();
    $(".t_main_forms").toggleClass("hide");
  });
  $("#field_byys_ration").change(function () {
    checkboxCountCheck();
    $(".t_byys_ration").toggleClass("hide");
  });
  $("#field_report_date_name").change(function () {
    checkboxCountCheck();
    $(".t_report_date_name").toggleClass("hide");
  });
  $("#field_report_opinion").change(function () {
    checkboxCountCheck();
    $(".t_report_opinion").toggleClass("hide");
  });
  $("#field_jzpg").change(function () {
    checkboxCountCheck();
    $(".t_jzpg").toggleClass("hide");
  });
  $("#field_latest_roe").change(function () {
    checkboxCountCheck();
    $(".t_latest_roe").toggleClass("hide");
  });
  $("#field_roe_tbzz").change(function () {
    checkboxCountCheck();
    $(".t_roe_tbzz").toggleClass("hide");
  });
  $("#field_roe_5y").change(function () {
    checkboxCountCheck();
    $(".t_roe_5y").toggleClass("hide");
  });
  $("#field_latest_eps").change(function () {
    checkboxCountCheck();
    $(".t_latest_eps").toggleClass("hide");
  });
  $("#field_eps_tbzz").change(function () {
    checkboxCountCheck();
    $(".t_eps_tbzz").toggleClass("hide");
  });
  $("#field_eps_5y").change(function () {
    checkboxCountCheck();
    $(".t_eps_5y").toggleClass("hide");
  });
  $("#field_total_income").change(function () {
    checkboxCountCheck();
    $(".t_total_income").toggleClass("hide");
  });
  $("#field_total_income_tbzz").change(function () {
    checkboxCountCheck();
    $(".t_total_income_tbzz").toggleClass("hide");
  });
  $("#field_total_income_5y").change(function () {
    checkboxCountCheck();
    $(".t_total_income_5y").toggleClass("hide");
  });
  $("#field_net_profit").change(function () {
    checkboxCountCheck();
    $(".t_net_profit").toggleClass("hide");
  });
  $("#field_net_profit_tbzz").change(function () {
    checkboxCountCheck();
    $(".t_net_profit_tbzz").toggleClass("hide");
  });
  $("#field_net_profit_5y").change(function () {
    checkboxCountCheck();
    $(".t_net_profit_5y").toggleClass("hide");
  });
  $("#field_zxgxl").change(function () {
    checkboxCountCheck();
    $(".t_zxgxl").toggleClass("hide");
  });
  $("#field_fina_report_date").change(function () {
    checkboxCountCheck();
    $(".t_fina_report_date").toggleClass("hide");
  });
  $("#field_fina_appoint_publish_date").change(function () {
    checkboxCountCheck();
    $(".t_fina_appoint_publish_date").toggleClass("hide");
  });
  $("#field_fina_actual_publish_date").change(function () {
    checkboxCountCheck();
    $(".t_fina_actual_publish_date").toggleClass("hide");
  });
  $("#field_total_market_cap").change(function () {
    checkboxCountCheck();
    $(".t_total_market_cap").toggleClass("hide");
  });
  $("#field_price").change(function () {
    checkboxCountCheck();
    $(".t_price").toggleClass("hide");
  });
  $("#field_right_price").change(function () {
    checkboxCountCheck();
    $(".t_right_price").toggleClass("hide");
  });
  $("#field_price_space").change(function () {
    checkboxCountCheck();
    $(".t_price_space").toggleClass("hide");
  });
  $("#field_hv").change(function () {
    checkboxCountCheck();
    $(".t_hv").toggleClass("hide");
  });
  $("#field_zxfzl").change(function () {
    checkboxCountCheck();
    $(".t_zxfzl").toggleClass("hide");
  });
  $("#field_fzldb").change(function () {
    checkboxCountCheck();
    $(".t_fzldb").toggleClass("hide");
  });
  $("#field_netprofit_growthrate_3_y").change(function () {
    checkboxCountCheck();
    $(".t_netprofit_growthrate_3_y").toggleClass("hide");
  });
  $("#field_income_growthrate_3_y").change(function () {
    checkboxCountCheck();
    $(".t_income_growthrate_3_y").toggleClass("hide");
  });
  $("#field_listing_yield_year").change(function () {
    checkboxCountCheck();
    $(".t_listing_yield_year").toggleClass("hide");
  });
  $("#field_listing_volatility_year").change(function () {
    checkboxCountCheck();
    $(".t_listing_volatility_year").toggleClass("hide");
  });
  $("#field_pe").change(function () {
    checkboxCountCheck();
    $(".t_pe").toggleClass("hide");
  });
  $("#field_peg").change(function () {
    checkboxCountCheck();
    $(".t_peg").toggleClass("hide");
  });
  $("#field_org_rating").change(function () {
    checkboxCountCheck();
    $(".t_org_rating").toggleClass("hide");
  });
  $("#field_profit_predict").change(function () {
    $(".t_profit_predict").toggleClass("hide");
  });
  $("#field_valuation_syl").change(function () {
    checkboxCountCheck();
    $(".t_valuation_syl").toggleClass("hide");
  });
  $("#field_valuation_sjl").change(function () {
    checkboxCountCheck();
    $(".t_valuation_sjl").toggleClass("hide");
  });
  $("#field_valuation_sxol").change(function () {
    checkboxCountCheck();
    $(".t_valuation_sxol").toggleClass("hide");
  });
  $("#field_valuation_sxnl").change(function () {
    checkboxCountCheck();
    $(".t_valuation_sxnl").toggleClass("hide");
  });
  $("#field_hyjzsp").change(function () {
    checkboxCountCheck();
    $(".t_hyjzsp").toggleClass("hide");
  });
  $("#field_ztzd").change(function () {
    checkboxCountCheck();
    $(".t_ztzd").toggleClass("hide");
  });
  $("#field_mll_5y").change(function () {
    checkboxCountCheck();
    $(".t_mll_5y").toggleClass("hide");
  });
  $("#field_jll_5y").change(function () {
    checkboxCountCheck();
    $(".t_jll_5y").toggleClass("hide");
  });
  $("#field_listing_date").change(function () {
    checkboxCountCheck();
    $(".t_listing_date").toggleClass("hide");
  });
  $("#field_netcash_operate").change(function () {
    checkboxCountCheck();
    $(".t_netcash_operate").toggleClass("hide");
  });
  $("#field_netcash_invest").change(function () {
    checkboxCountCheck();
    $(".t_netcash_invest").toggleClass("hide");
  });
  $("#field_netcash_finance").change(function () {
    checkboxCountCheck();
    $(".t_netcash_finance").toggleClass("hide");
  });
  $("#field_netcash_free").change(function () {
    checkboxCountCheck();
    $(".t_netcash_free").toggleClass("hide");
  });
});
