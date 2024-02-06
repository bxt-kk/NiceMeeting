function set_goto(btn, tag) {
    var navbar = document.querySelector(tag);
    var button = document.querySelector(btn);
    button.onclick = function() {
        window.scrollTo(0, 0);
    };
    document.addEventListener('scroll', function() {
        if (navbar.getBoundingClientRect().bottom < 0) {
            button.style.visibility = 'visible';
        } else {
            button.style.visibility = 'hidden';
        }
    });
}

function set_height(object, limit) {
    var U = document.querySelector(object);
    var D = document.querySelector(limit);
    const top = U.getBoundingClientRect().top;
    const height = D.getBoundingClientRect().height;
    U.style.minHeight = window.innerHeight - top - height + 'px';
}

function select_template(id) {
    var template = document.querySelector(id);
    var container = template.parentNode;
    container.removeChild(template);
    template.removeAttribute('id');
    return {tpl: template, ptr: container};
}

function extend_columns(data, columns, onclick) {
    for (var i = 0; i < columns.length; ++i) {
        var column = columns[i];
        var name = column.getAttribute('f-name');
        if (name != null) {
            var value = data[name];
            switch (column.tagName) {
                case 'IMG':
                    // column.src = `/img/${value}`;
                    column.src = `https://api.dicebear.com/7.x/lorelei/svg?seed=meeting`;
                    break;
                case 'TIME':
                    const time = new Date(value * 1000).toLocaleString();
                    column.setAttribute('datetime', time);
                    column.innerText = time;
                    break;
                case 'SUMMARY':
                    column.innerText = value.slice(0, 32);
                    break;
                default:
                    column.innerText = value;
                    break;
            }
        }
        var param = column.getAttribute('f-onclick');
        if (param != null && onclick != null) {
            const arg = data[param];
            column.onclick = function() {
                onclick(arg);
            };
            column.style.cursor = "pointer";
        }
        var cases = column.getAttribute('f-switch');
        if (cases != null && name != null) {
            column.innerText = cases.split(';')[column.innerText];
        }
    }
}

function extend_items(items, template, container, onclick=null) {
    items.forEach(item => {
        var row = template.cloneNode(true);
        extend_columns(item, row.children, onclick);
        container.append(row);
    });
}

function send_by_form(next, id="") {
    if (id === "") {
        var form = document.querySelector('form');
    } else {
        var form = document.getElementById(id);
    }
    var formData = new FormData(form);
    var jsonData = JSON.stringify(Object.fromEntries(formData));
    fetch(form.action, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: jsonData,
    }).then(
        resp => resp.json()
    ).then(data => {
        if (location.pathname == "/login.html") {
            next = `${next}?id=${data.id}`;
            localStorage.setItem('nm-user-id', `${data.id}`);
            localStorage.setItem('nm-user-level', `${data.level}`);
        }
        if (next !== "") location.href = next;
    });
    return false;
}

async function load_meeting(id) {
    const resp = await fetch(`/api/meeting/${id}`);
    return await resp.json();
}

async function load_meetings(page, size, category='') {
    const resp = await fetch(`/api/meetings/${size}/${page}/${category}`);
    return await resp.json();
}

async function load_lesson(id) {
    const resp = await fetch(`/api/lesson?id=${id}`);
    return await resp.json();
}

async function load_lessons(page, size) {
    const resp = await fetch(`/api/lessons?page=${page}&size=${size}`);
    return await resp.json();
}

async function load_contents(id) {
    const resp = await fetch(`/api/contents?id=${id}`);
    return await resp.json();
}

async function load_users(page, size) {
    const resp = await fetch(`/api/users?page=${page}&size=${size}`);
    return await resp.json();
}

async function load_user() {
    var id = (new URLSearchParams(location.search)).get('id') || '0';
    const resp = await fetch(`/api/user/${id}`);
    if (resp.status == 403) {
        location.href = '/login.html';
        return;
    }
    return await resp.json();
}

async function load_likes() {
    const resp = await fetch(`/api/likes`);
    return await resp.json();
}

async function in_likes(id) {
    const resp = await fetch(`/api/inlikes?course=${id}`);
    if (!resp.ok) return false;
    return await resp.json();
}

async function load_mimes(count, total) {
    const resp = await fetch(`/api/getmimes?count=${count}&total=${total}`);
    return await resp.json();
}

async function make_mimes(count, total) {
    const resp = await fetch(`/api/genmimes?count=${count}&total=${total}`);
    return await resp.text();
}

function set_loading(id) {
    var box = document.querySelector(id);
    if (box.children.length == 0) {
        var lds = document.createElement('DIV');
        lds.classList.add('lds-default');
        box.append(lds);
        for (let i = 0; i < 12; ++i) {
            lds.append(document.createElement('DIV'));
        }
    }
}

function del_loading(id) {
    var box = document.querySelector(id);
    if (box.children.length > 0) {
        box.removeChild(box.firstChild);
    }
}

function check_logout() {
    if ((new URLSearchParams(location.search)).get('from') == 'logout') {
        localStorage.removeItem('fl-login');
    }
}
