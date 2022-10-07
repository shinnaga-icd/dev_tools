var input_list = [];
var amount_maxlen = 8;

$(document).ready(function () {
  $('.form-control').
    focusin(function(evt) {
      $(this).addClass('input_current');
      if ($(this).attr('id') == "amount") {
        $(this).numberSeparator(false);
      }
    })
    .blur(function(evt) {
      $(this).removeClass('input_current');
      $(this).validateBlur();
    })
    .keyup(function(evt) {
      $(this).validateKeyUp();
    })
  ;
});

$(window).on('load',function(){
  if ($.trim($('#entryAlert').text())) {
    $('#entryAlert').removeClass('hidden');
  }
  $('.form-control').each(function(i, elem) {
    input_list.push( $(elem).attr('id') );
    var txt = $(elem).data('err');
    if (txt) {
      $(elem).validateResult(txt);
    }
    if ($(this).attr('id') == "amount") {
        $(this).numberSeparator(true);
    }
  });
});


$(document).on('click', '#btnConfirm', function(evt) {
  var first_err = null;
  if (!$('#entryAlert').hasClass('hidden')) {
      $('#entryAlert').addClass('hidden');
  }
  $('.form-control').each(function(i, elem) {
    if (!$(elem).validateBlur() && !first_err) {
      first_err = $(elem);
    }
  });
  //有効期限チェック
  if (!first_err) {
      var ykkMonth = "";
        var mObj = document.getElementsByName("W13D402.ykkMonth")[0];
        if (mObj) {
            for (var i = 0; i < mObj.length; i++) {
                if (mObj[i].selected) {
                    ykkMonth = mObj[i].value;
                }
            }
        }
        var ykkYear = "";
        var yObj = document.getElementsByName("W13D402.ykkYear")[0];
        if (yObj) {
            for (var i = 0; i < yObj.length; i++) {
                if (yObj[i].selected) {
                    ykkYear = yObj[i].value;
                }
            }
        }
        var ykk = ykkYear + ykkMonth;
        var curYm = document.getElementsByName("curYm")[0].value;

        if (ykk < curYm) {
            $('#entryAlert').text('有効期限切れのクレジットカードはご登録になれません。');
            $('#entryAlert').removeClass('hidden');
            var top = $('.alert').offset().top;
            $('html, body').animate({scrollTop:top});
            return;
        }
  }
  if (first_err) {
    var top = first_err.closest('.form-group').offset().top;
    $('html, body').animate({scrollTop:top});
  }
  else {
    // 精査チェックOKの場合
    // クレジットカード番号
    var first = document.getElementsByName("W13D402.firstIn")[0].value;
    var second = document.getElementsByName("W13D402.secondIn")[0].value;
    var third = document.getElementsByName("W13D402.thirdIn")[0].value;
    var fourth = document.getElementsByName("W13D402.fourthIn")[0].value;
    // 有効期限
    var ykkMonth = "";
    var mObj = document.getElementsByName("W13D402.ykkMonth")[0];
    if (mObj) {
        for (var i = 0; i < mObj.length; i++) {
            if (mObj[i].selected) {
                ykkMonth = mObj[i].value;
            }
        }
    }
    var ykkYear = "";
    var yObj = document.getElementsByName("W13D402.ykkYear")[0];
    if (yObj) {
        for (var i = 0; i < yObj.length; i++) {
            if (yObj[i].selected) {
                ykkYear = yObj[i].text;
            }
        }
    }

    // セキュリティコード
    var sctCod = document.getElementsByName("W13D402.sctCod")[0].value;
    // トークンAPIキー
    var tokenApiKey = document.getElementsByName("tokenApiKey")[0].value;
    // URL
    var url = document.getElementsByName('tokenApiUrl')[0].value;

    // トークンサーバとの通信処理
    var data = {};
    data.token_api_key = tokenApiKey;
    data.card_number = first + second + third + fourth;
    data.card_expire = ykkMonth + "/" + ykkYear;
    data.security_code = sctCod;
    data.lang = "ja";

    var xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.setRequestHeader('Accept', 'application/json');
    xhr.setRequestHeader('Content-Type', 'application/json; charset=utf-8');
    xhr.addEventListener('loadend', function () {
        if (xhr.status === 0) {
            $('#entryAlert').text('トークンサーバーとの接続に失敗しました');
            $('#entryAlert').removeClass('hidden');
            var top = $('.alert').offset().top;
            $('html, body').animate({scrollTop:top});
            return;
        }
        var response = JSON.parse(xhr.response);
        // 接続成功時
        if (xhr.status == 200) {
            // サーバー送信しない項目を初期化
            document.getElementsByName("W13D402.sctCod")[0].value = "";
            document.getElementsByName("W13D402.first")[0].value = first;
            document.getElementsByName("W13D402.second")[0].value = second;
            document.getElementsByName("W13D402.third")[0].value = "0000";
            document.getElementsByName("W13D402.fourth")[0].value = fourth;
            document.getElementsByName("W13D402.ksiKgk")[0].value = document.getElementsByName("W13D402.ksiKgkIn")[0].value.replace(/,/g, '');
            document.getElementsByName("W13D402.token")[0].value = response.token;
            document.getElementById('checkBtn').click();
        }
        else {
            $('#entryAlert').text(response.message);
            $('#entryAlert').removeClass('hidden');
            var top = $('.alert').offset().top;
            $('html, body').animate({scrollTop:top});
        }
    });
    xhr.send(JSON.stringify(data));
  }
});


$.fn.validateBlur = function() {

  var id = $(this).attr('id');
  var val = $(this).val();
  var maxlen = $(this).attr('maxlength');
  var err_msg = '';

  if (id == 'amount') {
    val = $(this).numberSeparator(false);
  }

  if ( $(this).attr('type') == 'text' || $(this).attr('type') == 'tel') {
    if (val.length == '0') {
      err_msg = '未入力です';
    }
    else if (!(/(^\d+$)/).test(val)) {
      err_msg = '半角数字で入力してください';
    }
    else {

      if (id == 'amount') {
        if (val.length > maxlen) {
          err_msg = maxlen + '桁以内で入力してください';
        }
        else if (val == '0') {
          err_msg = '0以上を入力してください';
        }

        $(this).numberSeparator(true);
      }
      else if (val.length != maxlen) {
        err_msg = maxlen + '桁入力してください';
      }

    }
    $(this).validateResult(err_msg);
  }

  return !(err_msg);
}

$.fn.validateKeyUp = function() {

  var id = $(this).attr('id');
  var val = $(this).val();
  var maxlen = $(this).attr('maxlength');
  var err_msg = '';

  if (id == 'amount') {
    val = $(this).numberSeparator(false);
  }

  if (($(this).attr('type') == 'text' || $(this).attr('type') == 'tel') && val.length > 0) {
    if (!(/(^\d+$)/).test(val)) {
      err_msg = '半角数字で入力してください';
    }
    else if (val.length > maxlen) {
      err_msg = maxlen + '桁以上入力されています';
    }
    $(this).validateResult(err_msg);

    if (val.length == maxlen) {
      if (id == 'first') {
        $("#second").focus();
      }
      else if (id == 'second') {
        $("#third").focus();
      }
      else if (id == 'third') {
        $("#fourth").focus();
      }
      else if (id == 'fourth') {
        $("#fourth").blur();
      }
    }

  }

  return !(err_msg);
}

$.fn.validateResult = function(err_message) {
  var id = $(this).attr('id');
  var elem_g = $(this).closest('.form-group');
  var elem_a = elem_g.children('.alert');
  var is_multiple = (elem_g.attr('id') == 'numberBlock');

  if (is_multiple) {
    var idx = $(this).parent().children('.form-control').index(this);
    var elem_a2 = elem_a.children('div').eq(idx);
  }

  if (err_message) {
    if (is_multiple) {
      elem_a2.text(err_message);
      elem_a2.removeClass('hidden');
    }
    else {
      elem_a.text(err_message);
    }
    elem_a.removeClass('hidden');
    $(this).addClass('input_error');
  }
  else {
    if (is_multiple) {
      elem_a2.empty(err_message);
      elem_a2.addClass('hidden');
      if (elem_a.children('.hidden').length == 4) {
        elem_a.addClass('hidden');
      }
    }
    else {
      elem_a.empty();
      elem_a.addClass('hidden');
    }
    $(this).removeClass('input_error');
  }
}

$.fn.numberSeparator = function($is_on) {
  var val = $(this).val();
  if ($is_on) {
    val = val.replace(/(\d)(?=(\d\d\d)+$)/g, '$1,');
    var cnt = (val.match(/,/g)||[]).length;
    $(this).attr("maxlength", (amount_maxlen + cnt));
  } else {
    val = val.replace(/,/g, '');
    $(this).attr("maxlength", amount_maxlen);
  }
  $(this).val(val);
  return val;
}
