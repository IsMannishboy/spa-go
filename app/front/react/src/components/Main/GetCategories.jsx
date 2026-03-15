import { useState,useEffect } from "react"
function Main(){
    const [categories,SetCategories] = useState([])
    useEffect(() => {
        fetch("/categories")
            .then(response => response.json())
            .then(data => {
                SetCategories(data);
            });
    }, []);

    return (
        <div>
            <h1>Categories</h1>
            {categories.length == 0 ?
             (<p>loading...</p>):
              <ul>
                {categories.map(category => (
                    <li key={category.id}>{category.name}</li>
                ))}
            </ul>
            }
            
        </div>
    )
}
export default Main;