async function set_player(id, lesson_id, autoplay=0) {
    const data = await load_lesson(lesson_id);
    switch (data.Status) {
        case 1:
            alert('缺銀子');
            return false;
        case 2:
            alert('請登錄');
            return false;
    }
    const source = data.Source;
    const player = document.querySelector(id);
    var frame = null;
    if (source.split('.').pop().toUpperCase() == 'M3U8') {
        frame = document.createElement('video')
        frame.controls = true;
        frame.autoplay = autoplay != 0;
        var hls = new Hls();
        hls.attachMedia(frame);
        hls.on(Hls.Events.MEDIA_ATTACHED, function () {
            hls.loadSource(source);
        });
    } else {
        var url = new URL(source);
        frame = document.createElement('iframe');
        frame.setAttribute('src', `https://www.ixigua.com/iframe${url.pathname}?autoplay=${autoplay}`);
        frame.setAttribute('frameborder', '0');
        frame.setAttribute('referrerpolicy', 'unsafe-url');
        frame.allowFullscreen = true;
    }
    if (frame != null) {
        player.childNodes.forEach(child => { child.remove(); });
        player.append(frame);
        frame.width = frame.parentElement.clientWidth;
        frame.height = Math.round(frame.width * 9 / 16);
        frame.style.maxWidth = '100%';
        window.addEventListener("resize", function() {
            frame.width = player.clientWidth;
            frame.height = Math.round(frame.width * 9 / 16);
        });
    }
    return true;
}
