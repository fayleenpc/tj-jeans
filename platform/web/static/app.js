let listProductHTML = document.querySelector('.listProduct');
let listCartHTML = document.querySelector('.listCart');
let iconCart = document.querySelector('.icon-cart');
let iconCartSpan = document.querySelector('.icon-cart span');
let body = document.querySelector('body');
let closeCart = document.querySelector('.close');
// let addCart = document.querySelector('.addCart')

let products = [];
let cart = [];



let checkoutCart = document.querySelector('.checkOut');


let register_submit = document.querySelector('.register-container .submit')
let login_submit = document.querySelector('.login-container .submit')


let cart_json = ``

let notifications = document.querySelector('.notifications');


let yippie = document.getElementsByClassName('yippie')

function getAuthCookie() {
    var cn = "Authorization=";
    var idx = document.cookie.indexOf(cn)
 
    if (idx != -1) {
        var end = document.cookie.indexOf(";", idx + 1);
        if (end == -1) end = document.cookie.length;
        return unescape(document.cookie.substring(idx + cn.length, end));
    } else {
        return "";
   }
}


function hasClass(el, className)
{
    if (el.classList)
        return el.classList.contains(className);
    return !!el.className.match(new RegExp('(\\s|^)' + className + '(\\s|$)'));
}

function addClass(el, className)
{
    if (el.classList)
        el.classList.add(className)
    else if (!hasClass(el, className))
        el.className += " " + className;
}

function removeClass(el, className)
{
    if (el.classList)
        el.classList.remove(className)
    else if (hasClass(el, className))
    {
        var reg = new RegExp('(\\s|^)' + className + '(\\s|$)');
        el.className = el.className.replace(reg, ' ');
    }
}
    
function createToast(type, icon, title, text){
    let newToast = document.createElement('div');
    newToast.innerHTML = `
        <div class="toast ${type}">
            <i class="${icon}"></i>
            <div class="content">
                <div class="title">${title}</div>
                <span>${text}</span>
            </div>
            <i class="fa-solid fa-xmark" onclick="(this.parentElement).remove()"></i>
        </div>`;
    notifications.appendChild(newToast);
    newToast.timeOut = setTimeout(
        ()=>newToast.remove(), 5000
    )
}



// success.onclick = function(){
//     let type = 'success';
//     let icon = 'fa-solid fa-circle-check';
//     let title = 'Success';
//     let text = 'This is a success toast.';
//     createToast(type, icon, title, text);
// }
// error.onclick = function(){
//     let type = 'error';
//     let icon = 'fa-solid fa-circle-exclamation';
//     let title = 'Error';
//     let text = 'This is a error toast.';
//     createToast(type, icon, title, text);
// }
// warning.onclick = function(){
//     let type = 'warning';
//     let icon = 'fa-solid fa-triangle-exclamation';
//     let title = 'Warning';
//     let text = 'This is a warning toast.';
//     createToast(type, icon, title, text);
// }
// info.onclick = function(){
//     let type = 'info';
//     let icon = 'fa-solid fa-circle-info';
//     let title = 'Info';
//     let text = 'This is a info toast.';
//     createToast(type, icon, title, text);
// }



function interact_login() {
    let show_user = document.querySelector('#nav-button')
    let show_user_body =  document.querySelector('#as_username')
    // let nav_menu = document.getElementById('navMenu')
    let nav_menu_layanan = document.querySelector('#layanan')

    show_user_body.innerHTML = `Login as ${localStorage.getItem('username')}`
    show_user.style.opacity = 1;
    nav_menu_layanan.style.opacity = 0;
}

function LogoutAS() {
    var logout_as = document.getElementById('logout')
    // var opacity_logout_as = getComputedStyle(logout_as).opacity;
    var nav_menu_layanan = document.querySelector('#layanan')
    var login_register_page = document.getElementById("login_register_page");
    var nav_layanan = document.getElementById("layanan");

    if (localStorage.getItem('username') != '') {
        login_register_page.style.display = "none";
        nav_layanan.style.display = "none";
        
        // logout_as.style.display = 'flex';
        logout_as.removeAttribute('disabled')
        logout_as.innerText = `Logout as ${localStorage.getItem('username')}`;
    } else {
        nav_menu_layanan.style.opacity = 1;
        // logout_as.style.display = 'none';
        // logout_as.removeAttribute('disabled')
        logout_as.innerText = `Logout as ${localStorage.getItem('username')}`;
    }
}

function Logout() {   
    // alert('logout')
    localStorage.clear()
    let type = 'info';
    let icon = 'fa-solid fa-circle-info';
    let title = 'Keluar Akun Berhasil';
    let text = 'Kamu akan keluar akun dalam beberapa detik dan dialihkan ke halaman utama.';
    createToast(type, icon, title, text);
    setTimeout(() => window.location = 'https://8383-103-147-8-149.ngrok-free.app', 5000)
}

login_submit.addEventListener('click', () => {
    // {
    //     "email": "me@me.com",
    //     "password": "asd"
    // }
    let username_login = document.querySelector('#login_username_email_input')
    let password_login = document.querySelector('#login_password_input')
    let start_login_json = `{ `;
    let body_login_json = ` "email": "${username_login.value}", "password": "${password_login.value}"`;
    let end_login_json = ` }`;
    let login_json = start_login_json + body_login_json + end_login_json
    let products_page = document.getElementById("products_page")
    
    // console.log(login_json)
    // alert('clicked login')
    fetch("/api/v1/login", {
        method: "POST",
        body: login_json,
        credentials: 'include',
        headers: {
          "Content-Type": "application/json; charset=UTF-8"
        }
    }).then(response => response.json())
    .then(data => {
        console.log(body_login_json)
        console.log(data)
        console.log('get cookie authorization header : ', getAuthCookie())
        if (data.error) {
            // alert(data.error)
            let type = 'error';
            let icon = 'fa-solid fa-circle-exclamation';
            let title = 'Masuk Akun Gagal';
            let text = 'Masukkan username/email dan password dengan benar, beralih ke halaman utama.';
            createToast(type, icon, title, text);
            // localStorage.clear()
            setTimeout(() => window.location = 'https://8383-103-147-8-149.ngrok-free.app', 5000)
            
            
        } 
        // if (data.admin && data.token) {
            // console.log("is admin ? ", data.admin)
            // console.log("token admin : ", data.token)
        // }
        
        if (data.token) {
            localStorage.setItem("username", username_login.value)
            localStorage.setItem('Authorization', data.token)
            // window.document.cookie = `Authorization : ${data.token}`
            products_page.style.display = 'block';
            addClass(document.getElementById("produk"), "active")
            // alert('success login')
            let type = 'success';
            let icon = 'fa-solid fa-circle-check';
            let title = 'Masuk Akun Berhasil';
            let text = 'Kamu berhasil login juga dapat berbelanja di menu produk untuk checkout sekaligus bayar.';
            createToast(type, icon, title, text);
            setTimeout(LogoutAS(), 5000)
            //problematic
            // interact_login()
        } 
        
        
    })
    // console.log(username_login.value)
    // console.log(password_login.value)
    
    
})

register_submit.addEventListener('click', () => {
    // {
    //     "email": "me@me.com",
    //     "password": "asd",
    //     "firstName": "tiago",
    //     "lastName": "user"
    // }
    let firstName_register = document.querySelector('#register_firstname_input')
    let lastName_register = document.querySelector('#register_lastname_input')
    let email_register = document.querySelector('#register_email_input')
    let password_register = document.querySelector('#register_password_input')
    let start_register_json = `{ `;
    let body_register_json = `"email": "${email_register.value}", "password": "${password_register.value}", "firstName": "${firstName_register.value}", "lastName": "${lastName_register.value}"`;
    let end_register_json = ` }`;
    let register_json = start_register_json + body_register_json + end_register_json
    // console.log(register_json)
    // alert('clicked register')

    fetch("/api/v1/register", {
        method: "POST",
        body: register_json,
        headers: {
          "Content-Type": "application/json; charset=UTF-8"
        }
    }).then(response => response.json())
    .then(data => {
        console.log(body_register_json)
        console.log(data)
        if (data.verify_url) {
            // alert('success register')
            let type = 'success';
            let icon = 'fa-solid fa-circle-check';
            let title = 'Daftar Akun Berhasil';
            let text = `Kamu telah berhasil mendaftarkan akun , mohon verifikasi dengan mengakses link berikut (dalam 5 detik), agar akun dapat digunakan : ${data.verify_url}`;
            createToast(type, icon, title, text);
            // setTimeout(() => window.location = `${data.verify_url}`, 5000)
            setTimeout(() => window.open(`${data.verify_url}`), 5000)
            // window.open(`${data.verify_url}`)
        }
        
        if (data != null && data.error) {
            // alert(data.error)
            let type = 'error';
            let icon = 'fa-solid fa-circle-exclamation';
            let title = 'Daftar Akun Gagal';
            let text = 'Kamu gagal melakukan daftar akun lalu akan beralih ke halaman utama.';
            createToast(type, icon, title, text);
            setTimeout(() => window.location = 'https://8383-103-147-8-149.ngrok-free.app', 5000)
        } 
        
    })
    
    // console.log(firstName_register.value)
    // console.log(lastName_register.value)
    // console.log(email_register.value)
    // console.log(password_register.value)
    
})



// {
//     "payment" : {
//       "payment_type" : "alfamart"
//     },
//     "customer": {
//       "name" : "john",
//       "email" : "foo@bar.com",
//       "phone_number" : ""
//     },
//     "items" : [{
//       "name" : "support podcast",
//       "category" : "podcast",
//       "merchant": "imregi.com",
//       "description": "donasi podcast imre",
//       "qty": 1,
//       "price": 10000,
//       "currency": "IDR"
//     }, 
//     {
//       "name" : "gk1020h12",
//       "category" : "podcast gk1020h12",
//       "merchant": "imregi.com gk1020h12",
//       "description": "donasi gk1020h12",
//       "qty": 1,
//       "price": 50000,
//       "currency": "IDR"
//     }]
//   }

// function sendEmail() {
//     Email.send({
//         Host: "smtp.yourisp.com",
//         Username: "username",
//         Password: "password",
//         To: 'recipient@example.com',
//         From: "sender@example.com",
//         Subject: "Test Email",
//         Body: "This is a test email sent using SMTP.js"
//     })
//     .then(function (message) {
//         alert("Mail sent successfully") // Alert message on successful email delivery
//     });
// }

iconCart.addEventListener('click', () => {
    body.classList.toggle('showCart');
})

closeCart.addEventListener('click', () => {
    body.classList.toggle('showCart');
})


checkoutCart.addEventListener('click', () => {
    // console.log("Authorization : ", localStorage.getItem("Authorization"))
    var messageForWhatsapp = ``
    var total = 0
    if (cart_json != ``) {
        fetch("/api/v1/cart/checkout", {
            method: "POST",
            body: cart_json,
            headers: {
              "Content-Type": "application/json; charset=UTF-8",
              "Authorization": localStorage.getItem("Authorization")
            }
        }).then(response => response.json())
        .then(data => {
            // console.log(data)
            if (data.error == 'permission denied') {
                // alert(data.error)
                let type = 'warning';
                let icon = 'fa-solid fa-triangle-exclamation';
                let title = 'Simpan Keranjang Belanja Gagal';
                let text = 'Kamu harus login terlebih dahulu untuk melakukan penyimpanan pembelian.';
                createToast(type, icon, title, text);
            } else {
                // alert('success checkout')
                // console.log(cart_json)
                // console.log(data)
                let type = 'success';
                let icon = 'fa-solid fa-circle-check';
                let title = 'Simpan Keranjang Belanja Berhasil';
                let text = 'Kamu telah menyimpan barang pembelian dan melakukan checkout/pembelian ke payment xendit/midtrans.';
                createToast(type, icon, title, text);
                let payment_invoices_json = `{ 
                "payment": {
                    "payment_type": "dana",
                    "amount": ${data.total_price}
                }, 
                "customer": {
                    "name": "TJ Jeans", 
                    "email": "${localStorage.getItem('username')}", 
                    "phone_number": "089505208391"
                }, 
                "items": ${JSON.stringify(data.items)}
                }`
                fetch("/api/v1/payment/invoices", {
                    method: "POST",
                    body: payment_invoices_json,
                    headers: {
                      "Content-Type": "application/json; charset=UTF-8",
                    }
                }).then(response => response.json())
                .then(data => {
                    console.log(data)
                    messageForWhatsapp += `Atas nama data penjual\n\nMerchant : ${data.billing_address.name}\nNo Telp : ${data.billing_address.phone_number}\n\nAtas nama data pembeli (isi data dengan lengkap)\n\nNama : \nEmail: ${localStorage.getItem('username')}\nNo Telp : \nAlamat Tujuan Pengiriman : \nBarang Yang Dibeli : \n`
                    for (i = 0; i < data.items.length; i++) {
                        messageForWhatsapp += `Produk -> ${data.items[i].name} -> Jumlah ${data.items[i].qty} -> Harga Satuan ${data.items[i].unit_price}\n`
                        total += data.items[i].qty * data.items[i].unit_price
                    }
                    messageForWhatsapp += `\nUntuk pembayaran totalnya seharga Rp ${total} (silahkan konfirmasi jika data sudah benar), lalu akses url berikut ${data.payment.redirect_url}\nMohon mengirimkan bukti transfer ketika sudah maka proses pengiriman dapat dilakukan, terima kasih.`

                    window.open('http://wa.me/089505208391?text='+ encodeURIComponent(messageForWhatsapp))
                })
            }
        })
    }
})



    const addDataToHTML = () => {
    // remove datas default from HTML

        // add new datas
        if(products.length > 0) // if has data
        {
            products.forEach(product => {
                let newProduct = document.createElement('div');
                newProduct.dataset.id = product.id;
                newProduct.classList.add('item');
                newProduct.innerHTML = 
                `<img src="/static/${product.image}" alt="${product.name}" widht="200px" height="200px">
                <h2>${product.name}</h2>
                <div class="price">${product.currency} ${product.price}</div>
                <div class="quantity">Qty ${product.qty}</div>
                <div class="merchant">@ ${product.merchant}</div>
                <div class="category">Category ${product.category}</div>
                <p>${product.description}</p>
                <button class="addCart" onclick="addClass(document.querySelector('.icon-cart'), 'yippie');setTimeout(() => removeClass(document.querySelector('.icon-cart'), 'yippie'),700)">Add To Cart</button>`;
                listProductHTML.appendChild(newProduct);
            });
        }
    }
    listProductHTML.addEventListener('click', (event) => {
        let positionClick = event.target;
        
        if(positionClick.classList.contains('addCart')){
            let id_product = positionClick.parentElement.dataset.id;
            
            addToCart(id_product);
        }
    })
const addToCart = (product_id) => {
    let positionThisProductInCart = cart.findIndex((value) => value.product_id == product_id);
    if(cart.length <= 0){
        cart = [{
            product_id: product_id,
            quantity: 1
        }];
    }else if(positionThisProductInCart < 0){
        cart.push({
            product_id: product_id,
            quantity: 1
        });
    }else{
        cart[positionThisProductInCart].quantity = cart[positionThisProductInCart].quantity + 1;
    }
    addCartToHTML();
    addCartToMemory();
}
const addCartToMemory = () => {
    localStorage.setItem('cart', JSON.stringify(cart));
}
const addCartToHTML = () => {
    listCartHTML.innerHTML = '';
    let totalQuantity = 0;

    let start_cart_json = `{ "items" : [  `;
    let body_cart_json = ``;
    let end_cart_json = `] }`;
    
    if(cart.length > 0){
        cart.forEach((item, index) => {
            totalQuantity = totalQuantity +  item.quantity;
            let newItem = document.createElement('div');
            newItem.classList.add('item');
            newItem.dataset.id = item.product_id;

            let positionProduct = products.findIndex((value) => value.id == item.product_id);
            let info = products[positionProduct];
            listCartHTML.appendChild(newItem);
            newItem.innerHTML = `
            <div class="image">
                    <img src="/static/${info.image}">
                </div>
                <div class="name">
                ${info.name}
                </div>
                <div class="totalPrice">IDR ${info.price * item.quantity}</div>
                <div class="quantity">
                    <span class="minus"><</span>
                    <span>${item.quantity}</span>
                    <span class="plus">></span>
                </div>
            `;
            body_cart_json += `{ "productID" : ${positionProduct+1}, "qty": ${item.quantity}}`;
            if (index != cart.length-1) {
                body_cart_json += `, `
            }
            
            
        })
        cart_json = start_cart_json + body_cart_json + end_cart_json
        // console.log(cart_json)
    }
    iconCartSpan.innerText = totalQuantity;
}

// {
//     "items": [
//       {
//         "productID": 1,
//         "quantity": 2
//       },
//       {
//         "productID": 2,
//         "quantity": 3
//       }
//     ]
// }

listCartHTML.addEventListener('click', (event) => {
    let positionClick = event.target;
    if(positionClick.classList.contains('minus') || positionClick.classList.contains('plus')){
        let product_id = positionClick.parentElement.parentElement.dataset.id;
        let type = 'minus';
        if(positionClick.classList.contains('plus')){
            type = 'plus';
        }
        changeQuantityCart(product_id, type);
    }
})
const changeQuantityCart = (product_id, type) => {
    let positionItemInCart = cart.findIndex((value) => value.product_id == product_id);
    if(positionItemInCart >= 0){
        let info = cart[positionItemInCart];
        switch (type) {
            case 'plus':
                cart[positionItemInCart].quantity = cart[positionItemInCart].quantity + 1;
                break;
        
            default:
                let changeQuantity = cart[positionItemInCart].quantity - 1;
                if (changeQuantity > 0) {
                    cart[positionItemInCart].quantity = changeQuantity;
                }else{
                    cart.splice(positionItemInCart, 1);
                }
                break;
        }
    }
    addCartToHTML();
    addCartToMemory();
}

// const initApp_Demo = () => {
//     // get data product
//     fetch('/static/products.json')
//     .then(response => response.json())
//     .then(data => {
//         products = data;
//         addDataToHTML();

//         // get data cart from memory
        // if(localStorage.getItem('cart')){
        //     cart = JSON.parse(localStorage.getItem('cart'));
        //     addCartToHTML();
        // }
//     })
// }
// initApp_Demo();
// fetch('/api/v1/product').then(response => response.json()).then(data => console.log(data))

const initApp = () => {
    // get data product
    fetch('/api/v1/product')
    .then(response => response.json())
    .then(data => {
        // console.log(data)
        products = data;
        addDataToHTML();

        // get data cart from memory
        if(localStorage.getItem('cart')){
            cart = JSON.parse(localStorage.getItem('cart'));
            addCartToHTML();
        }
        addCart = document.querySelector('.addCart')
        
    })
}


initApp();
