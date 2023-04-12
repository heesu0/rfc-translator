import React, { useState, useEffect } from "react";
import axios from "axios";

function App() {
  const [message, setMessage] = useState("");

  useEffect(() => {
    fetchMessage();
  }, []);

  const fetchMessage = async () => {
    try {
      const response = await axios.get("/api/rfc");
      setMessage(response.data);
    } catch (error) {
      console.error("Error fetching message:", error);
    }
  };

  return (
    <div>
      <h1>RFC 5246</h1>
      <pre>{message}</pre>
    </div>
  );
}

export default App;