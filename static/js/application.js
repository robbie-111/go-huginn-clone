// Huginn Application JavaScript Bundle
// CDN 라이브러리들 (base.templ에서 이미 로드됨):
//   - jQuery 2.x
//   - Bootstrap 3
//   - Select2 4.x
//   - Typeahead.js
//   - Spectrum

// Component files loaded here
// (These are loaded via separate <script> tags or bundled below)

// ============================================================
// Application initialization
// ============================================================
$(function () {
  // Initialize Select2 on all .select2-linked-tags elements
  if (typeof $.fn.select2 !== 'undefined') {
    $('.select2-linked-tags').select2({
      theme: 'bootstrap',
      templateResult: function(state) {
        if (!state.id) return state.text;
        var urlPrefix = $(state.element).parent().data('url-prefix');
        if (urlPrefix) {
          return $('<span><a href="' + urlPrefix + '/' + state.id + '" style="float:right;font-size:0.85em;" onclick="event.stopPropagation();" target="_blank">View</a>' + state.text + '</span>');
        }
        return state.text;
      }
    });
  }

  // Initialize Bootstrap popovers
  $('[data-toggle="popover"]').popover();

  // Initialize hover-help tooltips (Bootstrap popover)
  $('.hover-help').popover({
    trigger: 'hover',
    html: true,
    placement: 'right',
  });

  // Flash messages: auto-hide after 5 seconds
  var $flash = $('.flash');
  if ($flash.length) {
    setTimeout(function () {
      $flash.slideUp(500);
    }, 5000);
  }

  // Toggle memory display on agent show page
  $('#toggle-memory').on('click', function(e) {
    e.preventDefault();
    var $pre = $(this).siblings('.memory');
    if ($pre.hasClass('hidden')) {
      $pre.removeClass('hidden');
      $(this).text('Hide');
    } else {
      $pre.addClass('hidden');
      $(this).text('Show');
    }
  });

  // data-method support (for DELETE/PUT links without Rails UJS)
  $(document).on('click', 'a[data-method]', function(e) {
    var method = $(this).data('method');
    if (method && method.toUpperCase() !== 'GET') {
      e.preventDefault();
      var href = $(this).attr('href');
      var confirmMsg = $(this).data('confirm');
      if (confirmMsg && !confirm(confirmMsg)) return;
      var form = $('<form method="POST" style="display:none;"></form>').attr('action', href);
      form.append('<input type="hidden" name="_method" value="' + method.toUpperCase() + '"/>');
      form.append('<input type="hidden" name="csrf_token" value="mock-csrf-token"/>');
      $('body').append(form);
      form.submit();
    }
  });

  // data-confirm support for forms
  $(document).on('submit', 'form[data-confirm]', function(e) {
    var msg = $(this).data('confirm');
    if (!confirm(msg)) e.preventDefault();
  });
});
