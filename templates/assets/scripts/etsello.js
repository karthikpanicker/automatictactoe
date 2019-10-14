 $( document ).ready(function() {
    $('#boards button').on("click",function(){ 
        var boardId = $(this).attr("id");
        loadBoardLists(boardId,$(this));
     });
     if ($('#boards button').length > 0) {
         var selectedElem = $('#boards button')[0];
        if ($('#boards button.active').length > 0){
            selectedElem = $('#boards button.active')[0];
        }
        var boardId = $(selectedElem).attr('id');
        loadBoardLists(boardId,$(selectedElem));
     }

     $("#saveGTaskConfig").on("click",function(){
        var $this = $(this);
        buttonLoading($this);
        var gTaskListId = $('#googleLists button').filter('.active').attr('id');
        var transactionFilter = parseInt($('#gtask-radio-set input:radio:checked').val(),10);
        $.ajax({
            type: "POST",
            url: "api/users/1/gtasks-details",
            data: JSON.stringify({"listId": gTaskListId, "transactionFilter": transactionFilter}),
            success: function(data){
                setTimeout(function () {
                    $this.html($this.data('original-text'));
                }, 1000);
            }
        });
     });

     $("#saveTrelloConfig").on("click",function(){
        var $this = $(this);
        buttonLoading($this);
        var bId = $('#boards button').filter('.active').attr('id');
        var lId = $('#boardLists button').filter('.active').attr('id');
        var fields = []
        if ($('#listing_desc').prop("checked")){
            fields.push("listing_desc");
        }
        if ($('#listing_image').prop("checked")){
            fields.push("listing_image");
        }
        if ($('#listing_buy_profile').prop("checked")){
            fields.push("listing_buy_profile");
        }
        if ($('#listing_link').prop("checked")){
            fields.push("listing_link");
        }
        if ($('#listing_buyer_variations').prop("checked")){
            fields.push("listing_buyer_variations");
        }

        var transactionFilter = parseInt($('#trello-radio-set input:radio:checked').val(),10);
        $.ajax({
            type: "POST",
            url: "api/users/1/trello-details",
            data: JSON.stringify({"boardId": bId, "listId": lId,
                "fieldsToUse": fields, "transactionFilter": transactionFilter}),
            success: function(data){
                setTimeout(function () {
                    $this.html($this.data('original-text'));
                }, 1000);
            }
          });
    });
    $("#saveTodoistConfig").on("click",function(){
        var $this = $(this);
        buttonLoading($this);
        var todoistProjectId = parseInt($('#todoistProjects button').filter('.active').attr('id'),10);
        var transactionFilter = parseInt($('#todoist-radio-set input:radio:checked').val(),10);
        $.ajax({
            type: "POST",
            url: "api/users/1/todoist-details",
            data: JSON.stringify({"projectId": todoistProjectId, "transactionFilter": transactionFilter}),
            success: function(data){
                setTimeout(function () {
                    $this.html($this.data('original-text'));
                }, 1000);
            }
        });
     });
    var left  = ($(window).width()/2)-(500/2);
    var top   = ($(window).height()/2)-(600/2);
    $("#etsy-authorize").on("click",function(){
        window.open('/apps/etsy/authorize','EtsyAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });
    $("#trello-authorize").on("click",function(){
        window.open('/apps/trello/authorize','TrelloAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });

    $('.github-link').on("click",function(e){
        e.preventDefault(); 
        window.open('https://github.com/karthikpanicker/etsello','_blank');
    });

    $("#gtask-authorize").on("click",function(){
        window.open('/apps/gtask/authorize','GoogleAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });

    $("#todoist-authorize").on("click",function(){
        window.open('/apps/todoist/authorize','TodoistAuthorize', "width=500, height=600, top=" + top + ", left=" + left);
    });

    $('#gTaskConfigModal').on('show.bs.modal', function () {
        $('#googleLists button').remove();
        $('#spinner').show();
        $.get( "api/users/1/gtask-lists", function( gTasksLists ) {
            $.each(gTasksLists, function( index, list ) {
                var activeClass = (list.isSelected) ? ' active' : '' 
                $('#googleLists').append('<button id="' + list.id 
                + '" type="button" class="list-group-item' + activeClass +'">' + list.title + '</button>');
            });
            $('#spinner').hide();
            markActiveButton($('#googleLists button'),$(this));
        });
    });
    $('#todoistConfigModal').on('show.bs.modal', function () {
        $('#todoistProjects button').remove();
        $('#todoist-spinner').show();
        $.get( "api/users/1/todoist-projects", function( projects ) {
            $.each(projects, function( index, project ) {
                var activeClass = (project.isSelected) ? ' active' : '' 
                $('#todoistProjects').append('<button id="' + project.id 
                + '" type="button" class="list-group-item' + activeClass +'">' + project.name + '</button>');
            });
            $('#todoist-spinner').hide();
            markActiveButton($('#todoistProjects button'),$(this));
        });
    });
});

function loadBoardLists(boardId,selectedBoard) {
    $('#boards button').removeClass('active');
    selectedBoard.addClass('active');
    $.get( "api/users/1/trello-boards/" + boardId + "/lists", function( boardLists ) {
        $('#boardLists button').remove();
        $.each(boardLists, function( index, list ) {
            var activeClass = (list.isSelected) ? ' active' : '' 
            $('#boardLists').append('<button id="' + list.id 
            + '" type="button" class="list-group-item' + activeClass +'">' + list.name + '</button>');
        });
        markActiveButton($('#boardLists button'),$(this));
      }); 
}

function markActiveButton(buttonGroup,currentButton) {
    if (buttonGroup.length > 0) {
        // Bind click action on newly added items
        buttonGroup.on("click",function(){ 
            buttonGroup.removeClass('active');
            $(this).addClass('active');
        });
    }
}

function buttonLoading($this) {
    var loadingText = '<i class="fa fa-circle-o-notch fa-spin"></i> Saving...';
    if ($this.html() !== loadingText) {
        $this.data('original-text', $this.html());
        $this.html(loadingText);
    }
}