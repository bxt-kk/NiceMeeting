package page

type Page struct {
    Template string
    Data map[string]string
}

var pages = []Page{
    {
        Template: "index.html",
        Data: map[string]string{
            "title": "Home",
        },
    },
    {
        Template: "user.html",
        Data: map[string]string{
            "title": "User",
        },
    },
    {
        Template: "login.html",
        Data: map[string]string{
            "title": "Login",
        },
    },
    {
        Template: "signup.html",
        Data: map[string]string{
            "title": "Signup",
        },
    },
}
