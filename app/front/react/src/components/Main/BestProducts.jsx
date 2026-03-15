import { useState,useEffect } from "react";
function Producty(){
    const URL = "http://localhost/"
    const [Products,SetProducts] = useState([])
    
    useEffect(()=>{
        return async function fetchproducts() {
            const res = await fetch(`${URL}BestProducts`,{
                method:"GET"
                
            })
            if(!res.ok){
                console.log(res)
                return
            }
    const data = await res.json()
    SetProducts(data)        
    console.log(Products)
            
        }
    },[])
    return(
        <div>
            <div>
            <h1>Best Products</h1>
            {Products.length === 0 ? (
                <p>Loading...</p>
            ) : (
                <ul>
                    {Products.map(product => (
                        <li key={product.id}>
                            <img
                                src={`${URL}images/${product.image}`}
                                alt={product.name}
                                width="150"
                            />
                            {product.name}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    </div>
    )
}
export default Producty;