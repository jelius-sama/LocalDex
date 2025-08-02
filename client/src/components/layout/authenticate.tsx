import { useState } from 'react';

export default function AuthFlow({ onSuccess }: { onSuccess: () => void }) {
  const [step, setStep] = useState('login');     // 'login' or 'otp'
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [otp, setOtp] = useState('');
  const [message, setMessage] = useState('');

  const sendLogin = async () => {
    setMessage('');
    try {
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || res.statusText);
      }
      setStep('otp');
    } catch (err) {
      setMessage(err instanceof Error && err.message || 'Login error');
    }
  };

  const verifyOtp = async () => {
    setMessage('');
    try {
      const res = await fetch('/api/auth/verify_otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, otp })
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || res.statusText);
      }
      // cookie is set by server; notify parent
      onSuccess?.();
    } catch (err) {
      setMessage(err instanceof Error && err.message || 'OTP verification failed');
    }
  };

  return (
    <div style={{ maxWidth: 360, margin: '2em auto', fontFamily: 'sans-serif' }}>
      {step === 'login' && (
        <>
          <h2>Login</h2>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            style={{ width: '100%', padding: 8, marginBottom: 12 }}
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            style={{ width: '100%', padding: 8, marginBottom: 12 }}
          />
          <button onClick={sendLogin} style={{ width: '100%', padding: 10 }}>
            Send OTP
          </button>
          {message && <p style={{ color: 'red' }}>{message}</p>}
        </>
      )}

      {step === 'otp' && (
        <>
          <h2>Enter OTP</h2>
          <input
            type="text"
            placeholder="One-Time Password"
            value={otp}
            onChange={e => setOtp(e.target.value)}
            style={{ width: '100%', padding: 8, marginBottom: 12 }}
          />
          <button onClick={verifyOtp} style={{ width: '100%', padding: 10 }}>
            Verify & Login
          </button>
          {message && <p style={{ color: 'red' }}>{message}</p>}
          <p style={{ marginTop: 16, fontSize: 12, color: '#666' }}>
            Didn’t get an OTP?&nbsp;
            <button
              onClick={() => setStep('login')}
              style={{ background: 'none', border: 'none', color: '#06c', cursor: 'pointer' }}
            >
              Go back
            </button>
          </p>
        </>
      )}
    </div>
  );
}
