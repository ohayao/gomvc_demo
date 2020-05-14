$(function(){
    var gid='abc',uid="u10001",hb=0,ws=null;
    $('#connect').on('click',function(){
        gid=$('#room').val();
        uid=$('#uid').val();
        ws=new WebSocket("ws://localhost:11225/admin/websocket?gid="+gid+'&uid='+uid);
        ws.onopen=function(){
            hb=setInterval(function(){
                ws.send('sys_ping');
            },1e3*15);
            $('body input').hide();
            $('#content,#submit').fadeIn();

        }
        ws.onerror=function(e){
            clearInterval(hb);
        };
        ws.onclose=function(e){
            clearInterval(hb);
        };
        ws.onmessage=msg;
    });
    $('#submit').click(function(){
        var v=$('#content #input').html();
        if(v!=''){
            $('#content').append('<p style="color:red">你说：'+v+'</p>');
            ws.send(v);
            $('#content #input').html('');
        }
    });
    function msg(e){
        if(e.data!='sys_pong'){
            if(e.data.indexOf('系统消息')>-1){
                $('#content').append('<p style="color:#000;font-size:10px;">'+e.data+'</p>');
            }else $('#content').append('<p style="color:green">'+e.data+'</p>');
        }
    }
    document.onkeydown=function(e){
        if ((e.ctrlKey)&&(e.keyCode==13)){
            $('#submit').trigger('click');
        }
    }
});
