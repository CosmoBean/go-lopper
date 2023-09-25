import React, { useState } from 'react';

const Section =() => {
  const [login, setLogin] = useState(false);
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');

  const handleLogin = () => {
    setLogin(true);
  };

  const handleUrlChange = (e) => {
    setUrl(e.target.value);
  };

  const handleShortenUrl = () => {
    setShortUrl('http://short.ly/' + Math.random().toString(36).substr(2, 5));
  };

  return (
    <div className="h-screen bg-gray-200 flex flex-col justify-center items-center">
  <div className="absolute top-0 right-0 m-4">
    {!login && (
      <button onClick={handleLogin} className="p-2 bg-blue-500 text-white rounded">Login</button>
    )}
  </div>
  {login && (
    <form className="w-1/2">
      <input
        type="text"
        placeholder="Enter URL"
        value={url}
        onChange={handleUrlChange}
        className="w-full p-2 mb-2 border rounded"
      />
      <input
        type="text"
        placeholder="Enter Redirect URL"
        className="w-full p-2 mb-2 border rounded"
      />
      <div className="flex items-center mb-2">
        <input type="checkbox" className="mr-2"/>
        <label>Use custom shortened URL</label>
      </div>
      <button onClick={handleShortenUrl} className="p-2 bg-blue-500 text-white rounded">Shorten URL</button>
      {shortUrl && (
        <div className="mt-2">
          <span>Your shortened URL is: </span>
          <a href={shortUrl} target="_blank" rel="noopener noreferrer" className="text-blue-500">{shortUrl}</a>
        </div>
      )}
    </form>
  )}
</div>
);
};

export default Section;