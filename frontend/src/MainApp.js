import { useState } from "react";

function App() {
  const [url, setUrl] = useState("");
  const [customShortUrl, setCustomShortUrl] = useState(""); // Состояние для собственной короткой ссылки
  const [shortUrl, setShortUrl] = useState("");
  const [shortID, setShortID] = useState("");
  const [clicks, setClicks] = useState("");
  const [expireTime, setExpireTime] = useState("");
  const [creationTime, setCreationTime] = useState("");
  const [visitors, setVisitors] = useState([]);
  const [color, setColor] = useState("000000");
  const [copySuccess, setCopySuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const apiUrl = process.env.REACT_APP_API_URL;
  const [expireAtInput, setExpireAtInput] = useState("")


  if (!apiUrl) {
    console.log("api url не задан. проверь .env файл")
  }

  const handleShorten = async () => {
    setLoading(true);
    setShortUrl("");        // очистка перед новым запросом
    setShortID("");
    setClicks("");
    setVisitors([]);
    setCreationTime("");
    setExpireTime("");
    setLoading(true);

     try {
    const response = await fetch(`${apiUrl}/links`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ 
        original: url,
        shortened: customShortUrl || undefined,
        expire_at: expireAtInput
          ? new Date(expireAtInput + ":00")
          : undefined
      }),
    });

    if (response.ok) {
      const data = await response.json();
      console.log("otvet servera: ", data)
      setShortUrl(data.shortened);
      setShortID(data.id);
      setCreationTime(data.created_at);
      setExpireTime(data.expire_at);
    } else {
      alert("Failed to shorten the URL.");
    }
  } catch (error) {
    console.error("Error:", error);
    alert("An error occurred.");
  } finally {
    setLoading(false);
  }
};

  const handleClicks = async () => {
    try {
      const response = await fetch(`${apiUrl}/links/${shortID}`);
      if (response.ok) {
        const data = await response.json();
        setClicks(data.clicks);
      } else {
        alert("Failed to fetch clicks.");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("Error occurred while fetching clicks.");
    }
  };

  const handleVisitor = async () => {
    try {
      const response = await fetch(`${apiUrl}/links/${shortID}/visitors`);
      if (response.ok) {
        const data = await response.json();
        setVisitors(data.visitors);
      } else {
        alert("Failed to fetch visitors.");
      }
    } catch (error) {
      console.error("Error:", error);
      alert("Error fetching visitors.");
    }
  };

  const handleCopyClick = () => {
    navigator.clipboard.writeText(`${apiUrl}/${shortUrl}`).then(() => {
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000); // Скрыть уведомление через 2 секунды
    });
  };

  return (
    <div className="bg-gradient-to-br from-purple-300 to-white text-purple-900 min-h-screen flex items-center justify-center px-4 py-10">
      <div className="flex flex-row items-start space-x-6 p-8 border border-purple-300 rounded-xl shadow-lg bg-white w-full max-w-6xl">
        {/* Left panel */}
        <div className="flex flex-col space-y-6 w-2/3">
          <h1 className="text-3xl font-bold">URL Shortener</h1>

          <input
            type="text"
            placeholder="Введите URL"
            className="p-2 bg-purple-50 text-purple-900 border border-purple-500 rounded focus:outline-none focus:ring-2 focus:ring-purple-700"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
          />

          {/* 🎨 Color Picker */}
          <div className="flex flex-col items-start space-y-2">
            <label htmlFor="qrColor" className="text-sm font-semibold text-purple-800">
              Цвет QR-кода
            </label>
            <div className="flex items-center space-x-3">
              <input
                type="color"
                id="qrColor"
                value={`#${color}`}
                onChange={(e) => setColor(e.target.value.replace("#", ""))}
                className="w-12 h-10 rounded-md border-2 border-purple-400 cursor-pointer shadow-md hover:shadow-lg transition"
                title="Выберите цвет QR-кода"
              />
              <span className="text-sm text-gray-600">
                {`#${color.toUpperCase()}`}
              </span>
              <div
                className="w-6 h-6 rounded-full border border-gray-300"
                style={{ backgroundColor: `#${color}` }}
              />
            </div>
            <p className="text-xs text-gray-500">
              Цвет применяется при генерации QR-кода
            </p>
          </div>

          <div className="flex flex-col space-y-2">
            <input
              type="text"
              placeholder="Введите свою короткую ссылку (опционально)"
              className="p-2 bg-purple-50 text-purple-900 border border-purple-500 rounded focus:outline-none focus:ring-2 focus:ring-purple-700"
              value={customShortUrl}
              onChange={(e) => setCustomShortUrl(e.target.value)}
            />
            <p className="text-xs text-gray-500">
              Оставьте это поле пустым, чтобы система создала случайное имя.
            </p>
          </div>

          <div className="flex flex-col space-y-2">
  <label htmlFor="expireAt" className="text-sm font-semibold text-purple-800">
    Время жизни ссылки (опционально)
  </label>
  <input
    type="datetime-local"
    id="expireAt"
    value={expireAtInput}
    onChange={(e) => setExpireAtInput(e.target.value)}
    className="p-2 bg-purple-50 text-purple-900 border border-purple-500 rounded focus:outline-none focus:ring-2 focus:ring-purple-700"
  />
  <p className="text-xs text-gray-500">
    Если не указано — ссылка будет действовать 24 часа по умолчанию.
  </p>
</div>


          <button
            onClick={handleShorten}
            className="py-2 bg-purple-700 text-white rounded hover:bg-purple-800 transition"
            disabled={loading}
          >
            {loading ? "Сокращение..." : "Сократить"}
          </button>


          {shortUrl && (
            <div>
              <p className="text-sm">Ваша ссылка:</p>
              <div className="flex items-center gap-2 text-center">
                <a
                  href={`${apiUrl}/${shortUrl}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-purple-700 underline break-all"
                >
                  {`${apiUrl}/${shortUrl}`}
                </a>
                <button
                  onClick={handleCopyClick}
                  className="text-sm text-purple-500 hover:underline"
                >
                  {copySuccess ? "Скопировано!" : "Копировать"}
                </button>
              </div>
              <div className="mt-3 text-sm text-gray-700">
                <p>Создано: {creationTime}</p>
                <p>Истекает: {expireTime} </p>
              </div>
            </div>
          )}

          <div className="pt-4">
            
            <button
              onClick={handleClicks}
              className="w-full mt-2 py-2 bg-purple-700 text-white rounded hover:bg-purple-800 transition"
            >
              Кол-во переходов
            </button>
            <p className="mt-2 text-sm">Переходов: {clicks}</p>
          </div>

          <div>
            <button
              onClick={handleVisitor}
              className="w-full mt-4 p-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition"
            >
              Показать посетителей
            </button>

            {visitors.length > 0 && (
              <div className="mt-4 w-full text-left max-h-60 overflow-y-auto">
                <h3 className="text-lg font-semibold mb-2">Последние посетители:</h3>
                {visitors.map((visitor) => (
                  <div
                    key={visitor.id}
                    className="text-sm border-b border-purple-200 py-1"
                  >
                    <strong>ID #{visitor.ID}</strong>: {visitor.user_agent} — {visitor.user_IP}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Right panel (QR) */}
        <div className="w-1/3 flex flex-col items-center">
          <h2 className="text-xl font-semibold mb-4">QR-Код</h2>
          {shortUrl ? (
            <img
              src={`${apiUrl}/qr/view?url=${apiUrl}/${shortUrl}&color=${color}`}
              alt="QR Code"
              className="w-64 h-64 border rounded-lg shadow-md"
            />
          ) : (
            <p className="text-sm text-gray-500 text-center">
              Сначала сократите ссылку, чтобы получить QR
            </p>
          )}
          {shortUrl && (
            <a
              href={`${apiUrl}/qr/download?url=${apiUrl}/${shortUrl}&color=${color}`}
              className="mt-4 px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition text-sm"
            >
              Скачать QR-код
            </a>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
