function new_menu(id, items=[]) {
    var node = document.querySelector(id);
    var link = document.createElement('A');
    link.classList.add('menu-logo');
    link.setAttribute('href', '/');
    node.append(link);
    var logo = document.createElement('IMG');
    logo.setAttribute('src', '/logo/Golang.png');
    link.append(logo);
    var head = document.createElement('A');
    head.classList.add('menu-title');
    head.setAttribute('href', '/');
    head.innerText = 'funnylesson';
    node.append(head);
    const path = location.pathname;
    if (!['/login.html', '/signup.html'].includes(path)) {
        var item = document.createElement('A');
        node.append(item);
        item.classList.add('menu-item');
        if (localStorage.getItem('fl-login') != 'true') {
            item.setAttribute('href', '/login.html');
            item.innerText = '登錄|註冊';
        } else if (path == '/user.html'){
            item.setAttribute('href', '/api/logout');
            item.innerText = '登出';
        } else {
            item.setAttribute('href', '/user.html');
            item.innerText = '個人主頁';
        }
    }
    items.forEach(item => {
        var a = document.createElement('A');
        node.append(a);
        a.classList.add('menu-item');
        a.setAttribute('href', item.href);
        a.innerText = item.text;
    });
}

async function fl_index(categories) {
    new_menu('#menu', [{href: '/courses.html', text: '課程'}]);
    set_goto('#goto', 'nav');

    var classes = document.querySelector('#classes');
    var tags = classes.querySelector('.tags');
    categories.forEach((category, i) => {
        var tag = document.createElement('a');
        tag.innerText = category;
        tag.href = '/courses.html';
        if (i > 0) tag.href += `?category=${category}`;
        tags.append(tag);
    });

    var tpl = select_template('#course').tpl;
    var nodes = select_template('#courses');
    for (let i = 1; i < categories.length; ++i) {
        var section = nodes.tpl.cloneNode(true)
        nodes.ptr.append(section);
        const category = categories[i];
        const data = await load_courses(0, 12, category);
        section.querySelector('header').innerText = category;
        var ptr = section.querySelector('.gallery');
        extend_items(data.Courses, tpl, ptr, function(id) {
            location.href = `/course.html?id=${id}`;
        });
    }
}

async function fl_courses(categories) {
    const category = (new URLSearchParams(location.search)).get('category') || '';
    var nodes = select_template('#course');
    window.PAGE = 0;

    async function load_more() {
        set_loading('#loading');
        var data = await load_courses(window.PAGE++, 4, category);
        del_loading('#loading');
        extend_items(data.Courses, nodes.tpl, nodes.ptr, function(id) {
            location.href = `/course.html?id=${id}`;
        });
        return data.Total > window.PAGE;
    }

    new_menu('#menu', [{href: '/', text: '首頁'}]);
    set_goto('#goto', 'nav');
    set_height('main', 'footer');

    var classes = document.querySelector('#classes');
    var header = classes.querySelector('header')
    header.innerText = category || header.innerText;
    var tags = classes.querySelector('.tags');
    categories.forEach((item, i) => {
        var tag = document.createElement('a');
        tags.append(tag);
        tag.innerText = item;
        tag.href = '/courses.html';
        if (i > 0) tag.href += `?category=${item}`;
        if (item == category || (i == 0 && category == '')) {
            tag.classList.add('selected');
            tag.removeAttribute('href');
        }
    });

    var footer = document.querySelector("footer");
    while (footer.getBoundingClientRect().top < window.innerHeight) {
        if (!(await load_more())) {
            break;
        }
    }

    window.onscroll = async function() {
        var tag = document.querySelector('footer');
        var top = tag.getBoundingClientRect().top;
        if (top < window.innerHeight) {
            const onscroll = window.onscroll;
            window.onscroll = null;
            if (await load_more()) {
                window.onscroll = onscroll;
            }
        }
    };
}

async function fl_course() {
    new_menu('#menu', [{href: '/courses.html', text: '課程'}]);
    set_goto('#goto', '#player');
    set_height('main', 'footer');

    const id = (new URLSearchParams(location.search)).get('id');
    const course = await load_course(id);
    var columns = document.querySelector('#course').children;
    extend_columns(course, columns, function(name) {
        location.href = `/courses.html?category=${name}`;
    });
    var isfollowd = await in_likes(course.Id);
    var followbtn = document.querySelector('#follow');
    followbtn.innerText = isfollowd ? '取關' : '關注';
    followbtn.onclick = async function() {
        var url = `/api/unfollow?course=${course.Id}`;
        if (!isfollowd) {
            url = `/api/follow?course=${course.Id}`;
        }
        const resp = await fetch(url);
        if (resp.ok) {
            isfollowd = !isfollowd;
            followbtn.innerText = isfollowd ? '取關' : '關注';
        } else if (resp.status == 403) {
            location.href = `/login.html`;
        }
    };
    if (course.Id == 0) return;

    set_loading('#loading');
    const data = await load_contents(course.Id);
    del_loading('#loading');
    var index = document.querySelector("#index");
    data.forEach((lesson, id) => {
        var item = document.createElement('A');
        index.append(item);
        item.innerText = lesson.Id;
        const sup = document.createElement('sup');
        const small = document.createElement('small');
        sup.append(small);
        switch (lesson.Status) {
            case 1:
                small.innerText = '付費';
                item.append(sup)
                break;
            case 2:
                small.innerText = '登錄';
                item.append(sup)
                break;
        }
        item.onclick = async function() {
            const ret = await set_player('#player', lesson.Id, 1);
            if (!ret) return;
            var columns = document.querySelector('#lesson').children;
            extend_columns(lesson, columns, null);
        };
    })
    var nodes = select_template('#lesson-item');
    extend_items(data, nodes.tpl, nodes.ptr, async function(id) {
        var lesson = data[id - 1];
        const ret = await set_player('#player', lesson.Id, 1);
        if (!ret) return;
        var columns = document.querySelector('#lesson').children;
        extend_columns(lesson, columns, null);
    });
    if (data.length > 0) {
        set_player('#player', data[0].Id);
        var columns = document.querySelector('#lesson').children;
        extend_columns(data[0], columns, null);
    }
}

async function fl_lessons() {
    var nodes = select_template('#lesson');
    window.PAGE = 0;

    async function load_more() {
        set_loading('#loading');
        var data = await load_lessons(window.PAGE++, 4);
        del_loading('#loading');
        extend_items(data.Lessons, nodes.tpl, nodes.ptr, function(id) {
            location.href = `/lesson.html?id=${id}`;
        });
        return data.Total > window.PAGE;
    }

    new_menu('#menu', [{href: '/courses.html', text: '課程'}]);
    set_goto('#goto', 'nav');
    set_height('main', 'footer');

    var footer = document.querySelector("footer");
    while (footer.getBoundingClientRect().top < window.innerHeight) {
        if (!(await load_more())) {
            break;
        }
    }

    window.onscroll = async function() {
        var tag = document.querySelector('footer');
        var top = tag.getBoundingClientRect().top;
        if (top < window.innerHeight) {
            const onscroll = window.onscroll;
            window.onscroll = null;
            if (await load_more()) {
                window.onscroll = onscroll;
            }
        }
    };
}

function fl_login() {
    new_menu('#menu', [{href: '/courses.html', text: '課程'}, {href: '/', text: '首頁'}]);
    check_logout();
}

function fl_signup() {
    new_menu('#menu', [{href: '/courses.html', text: '課程'}, {href: '/', text: '首頁'}]);
}

function fl_mining() {
    new_menu('#menu', [{href: '/courses.html', text: '課程'}]);
}

async function fl_user() {
    new_menu('#menu', [{href: '/courses.html', text: '課程'}]);
    set_height('main', 'footer');
    set_goto('#goto', 'nav');

    const user = await load_user();
    var columns = document.querySelector('#user-profile').children;
    extend_columns(user, columns, null);
    var src = `https://avatars.dicebear.com/api/big-smile/${user.Username}.svg`;
    document.querySelector('#user-picture').src = src;
    if (user.Level > 1) {
        const elem = document.querySelector('#level');
        const link = document.createElement('a');
        elem.append(link);
        link.href = '/worker/index.html';
        link.innerText = '#點擊進入工作頁';
        link.style.marginLeft = '0.1rem';
        link.style.fontSize = 'small';
    }

    set_loading('#loading');
    var data = await load_likes();
    del_loading('#loading');
    var nodes = select_template('#course');
    extend_items(data, nodes.tpl, nodes.ptr, function(id) {
        location.href = `/course.html?id=${id}`;
    });

    window.onscroll = function() {
        var navbar = document.querySelector('nav');
        var bottom = navbar.getBoundingClientRect().bottom;
        if (bottom < 0) {
            document.querySelector('#goto').style.visibility = 'visible';
        } else {
            document.querySelector('#goto').style.visibility = 'hidden';
        }
    };
}
