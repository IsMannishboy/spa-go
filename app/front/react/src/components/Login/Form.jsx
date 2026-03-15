import { useState,useEffect } from "react";

 function Form() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

     useEffect(() => {
            URL = "http://localhost:80/"

        fetch(`${URL}auth/csrf`, {
            credentials: "include",
        })
            .then((res)=>{
            const token = res.headers.get("CSRF");
            console.log("csrf fetched:", token);
            })
            
    }, []); 
    function Submit(e) {
        URL = "http://localhost:80/"

        e.preventDefault();

        const data = {
            username: username,
            password: password,
        };  
        console.log(data)
        fetch(`${URL}auth/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
               
            },
            credentials: "include",
            body: JSON.stringify(data),
        })
            .then((res)=>{
                if(res.ok){
                    console.log("loged in sucessfull")
                                    window.location.href = `${URL}main`

                }
            })
            
            .catch(err => console.error(err));
    }

    return (
        <form onSubmit={Submit}>
            <h1>Login</h1>

            <p>Username</p>
            <input
                name="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />

            <p>Password</p>
            <input
                type="password"
                name="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />

            <button type="submit">Login</button>
        </form>
    );
}

export default Form;
