import './app.css'
import { useState, useEffect } from 'preact/hooks'

export function App() {
  const [items, setItems] = useState([])
  const [itemId, setItemId] = useState('')
  const [itemName, setItemName] = useState('')
  const [itemQuantity, setitemQuantity] = useState('')
  const [message, setMessage] = useState('')

  const API_URL = 'http://localhost:8080' // ← Passe das an deinen Server an

  // GET /items
  const fetchAllItems = async () => {
    const res = await fetch(`${API_URL}/items`)
    const data = await res.json()
    console.log(data);
    
    setItems(data)
  }

  // GET /items/:itemId
  const fetchItemById = async () => {
    const res = await fetch(`${API_URL}/items/${itemId}`)
    const data = await res.json()
    setItems([data]) // Einzelnes Item als Array anzeigen
  }

  // POST /items
  const createOrUpdateItem = async () => {
    const res = await fetch(`${API_URL}/items`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: Number(itemId), name: itemName, quantity: Number(itemQuantity)}),
    })
    const data = await res.json()
    setMessage(data.message || 'Item erstellt/aktualisiert')
    fetchAllItems()
  }

  // PUT /items/:itemId
  const updateItem = async () => {
    const res = await fetch(`${API_URL}/items/${itemId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: Number(itemId), name: itemName, quantity: Number(itemQuantity) }),
    })
    const data = await res.json()
    setMessage(data.message || 'Item aktualisiert')
    fetchAllItems()
  }

  // DELETE /items/:itemId
  const deleteItem = async () => {
    const res = await fetch(`${API_URL}/items/${itemId}`, {
      method: 'DELETE',
    })
    const data = await res.json()
    setMessage(data.message || 'Item gelöscht')
    fetchAllItems()
  }

  useEffect(() => {
    fetchAllItems()
  }, [])

  return (
    <div class="app">
      <h1>Item Manager</h1>
      <div class="controls">
        <input
          type="text"
          placeholder="Item ID"
          value={itemId}
          onInput={(e) => setItemId(e.target.value)}
        />
        <input
          type="text"
          placeholder="Item Name"
          value={itemName}
          onInput={(e) => setItemName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Item Quanitiy"
          value={itemQuantity}
          onInput={(e) => setitemQuantity(e.target.value)}
        />
        <button onClick={createOrUpdateItem}>Create/Update (POST)</button>
        <button onClick={updateItem}>Update by ID (PUT)</button>
        <button onClick={deleteItem}>Delete by ID (DELETE)</button>
        <button onClick={fetchItemById}>Get by ID (GET)</button>
        <button onClick={fetchAllItems}>Get All Items</button>
      </div>

      {message && <p class="message">{message}</p>}

      <div class="items">
        <h2>Items</h2>
        {items.length === 0 && <p>No items found.</p>}
        {items.map((item) => (
          <div class="item" key={item.id}>
            <strong>ID:</strong> {item.id}<br />
            <strong>Name:</strong> {item.name}<br />
            <strong>Quantity:</strong> {item.quantity}<br />
          </div>
        ))}
      </div>
    </div>
  )
}