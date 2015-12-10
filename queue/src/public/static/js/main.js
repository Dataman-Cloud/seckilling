/**
 * Created by yaoyun on 15/12/10.
 */
function Flush() {
    $.ajax({
        type: "GET",
        url: "/v1/events/1",
        data: {},
        success: function(data){
            console.log(data["data"])

        }
    });
}

$(document).ready(function(){

});

// Common functions
function pad(number, length) {
    var str = '' + number;
    while (str.length < length) {str = '0' + str;}
    return str;
}

function formatTime(time) {
    var min = parseInt(time / 6000),
        sec = parseInt(time / 100) - (min * 60),
        hundredths = pad(time - (sec * 100) - (min * 6000), 2);
    return (min > 0 ? pad(min, 2) : "00") + ":" + pad(sec, 2) + ":" + hundredths;
}

var Example2 = new (function() {
    var $countdown,
        incrementTime = 70,
        currentTime = 3000,
        updateTimer = function() {
            $countdown.html("倒计时: " + formatTime(currentTime));
            if (currentTime == 0) {
                Example2.Timer.stop();
                timerComplete();
                Example2.resetCountdown();
                return;
            }
            currentTime -= incrementTime / 10;
            if (currentTime < 0) currentTime = 0;
        },
        timerComplete = function() {
            $("#btWait").hide()
            $("#btOver").hide()
            $("#btBuy").show()
            alert('活动开始');
        },
        init = function() {
            $countdown = $('#countdown');
            Example2.Timer = $.timer(updateTimer, incrementTime, true);
        };
    this.resetCountdown = function() {
        var newTime = 30 * 100;
        if (newTime > 0) {currentTime = newTime;}
        this.Timer.stop().once();
    };
    $(init);
});
