import { NextResponse } from 'next/server';
import { OAuth2Client } from 'google-auth-library';

// Inisialisasi OAuth2Client dengan Client ID kita dari .env.local
const client = new OAuth2Client(process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID);

export async function POST(request: Request) {
  try {
    // Ambil token dari body request yang dikirim frontend
    const body = await request.json();
    const { token } = body;

    if (!token) {
      return NextResponse.json(
        { error: 'Token tidak disediakan.' },
        { status: 400 }
      );
    }

    // Verifikasi token menggunakan OAuth2Client
    const ticket = await client.verifyIdToken({
      idToken: token,
      audience: process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID, 
      // Penting: Audience harus sama persis dengan Client ID aplikasi.
    });

    // Ambil payload/data user dari token yang sudah diverifikasi
    const payload = ticket.getPayload();

    if (!payload) {
      return NextResponse.json(
        { error: 'Gagal mendapatkan data payload dari token.' },
        { status: 401 }
      );
    }

    // Mengembalikan data sederhana sebagai kesuksesan verifikasi Google Token
    return NextResponse.json(
      { 
        message: 'Login berhasil diverifikasi',
        user: {
          email: payload.email,
          name: payload.name,
          picture: payload.picture,
          googleId: payload.sub // ID Unik dari Google
        },
        // Anda juga bisa men-generate JWT token (contoh aplikasi full Next.js)
        token: token // Mengirim balik token atau bisa ditukar dengan sesi internal Backend 
      },
      { status: 200 }
    );

  } catch (error) {
    console.error('Error memverifikasi token Google:', error);
    // Jika verifikasi gagal (token expired/invalid)
    return NextResponse.json(
      { error: 'Unauthorized: Token invalid atau tidak bisa diverifikasi.' },
      { status: 401 }
    );
  }
}
