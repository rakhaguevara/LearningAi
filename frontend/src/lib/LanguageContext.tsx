'use client';

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';

type Language = 'en' | 'id';

interface LanguageContextType {
    language: Language;
    setLanguage: (lang: Language) => void;
    t: (key: string) => string;
}

const translations = {
    en: {
        // Appearance
        'appearance.title': 'Appearance',
        'theme.label': 'Theme',
        'theme.description': 'Switch between dark and light mode',
        'theme.dark': '🌙 Dark',
        'theme.light': '☀️ Light',
        'language.label': 'Language',
        'language.description': 'Interface language',

        // AI Preferences
        'aiPreferences.title': 'AI Preferences',
        'ai.autoSummarize': 'Auto-summarize sources',
        'ai.autoSummarize.desc': 'Generate a summary when sources are added',
        'ai.adaptiveTone': 'Adaptive tone',
        'ai.adaptiveTone.desc': 'Adjust explanation style based on behavior signals',
        'ai.streamResponses': 'Stream responses',
        'ai.streamResponses.desc': 'Show AI responses as they\'re generated',

        // Pomodoro
        'pomodoro.title': 'Pomodoro',
        'pomodoro.soundAlerts': 'Sound alerts',
        'pomodoro.soundAlerts.desc': 'Play a sound when timer ends',
        'pomodoro.autoStart': 'Auto-start breaks',
        'pomodoro.autoStart.desc': 'Automatically start break timer',
        'pomodoro.notifications': 'Notifications',
        'pomodoro.notifications.desc': 'Browser notifications for session end',

        // Danger Zone
        'dangerZone.title': 'Danger Zone',
        'dangerZone.deleteAccount': 'Delete Account',
        'dangerZone.deleteAccount.desc': 'Permanently remove all data',
        'dangerZone.delete': 'Delete',
    },
    id: {
        // Appearance
        'appearance.title': 'Tampilan',
        'theme.label': 'Tema',
        'theme.description': 'Beralih antara mode gelap dan terang',
        'theme.dark': '🌙 Gelap',
        'theme.light': '☀️ Terang',
        'language.label': 'Bahasa',
        'language.description': 'Bahasa antarmuka',

        // AI Preferences
        'aiPreferences.title': 'Preferensi AI',
        'ai.autoSummarize': 'Ringkas sumber otomatis',
        'ai.autoSummarize.desc': 'Buat ringkasan saat sumber ditambahkan',
        'ai.adaptiveTone': 'Nada adaptif',
        'ai.adaptiveTone.desc': 'Sesuaikan gaya penjelasan berdasarkan perilaku',
        'ai.streamResponses': 'Respons streaming',
        'ai.streamResponses.desc': 'Tampilkan respons AI saat dibuat',

        // Pomodoro
        'pomodoro.title': 'Pomodoro',
        'pomodoro.soundAlerts': 'Peringatan suara',
        'pomodoro.soundAlerts.desc': 'Putar suara saat timer berakhir',
        'pomodoro.autoStart': 'Mulai istirahat otomatis',
        'pomodoro.autoStart.desc': 'Mulai timer istirahat secara otomatis',
        'pomodoro.notifications': 'Notifikasi',
        'pomodoro.notifications.desc': 'Notifikasi browser saat sesi berakhir',

        // Danger Zone
        'dangerZone.title': 'Zona Berbahaya',
        'dangerZone.deleteAccount': 'Hapus Akun',
        'dangerZone.deleteAccount.desc': 'Hapus semua data secara permanen',
        'dangerZone.delete': 'Hapus',
    },
};

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export function LanguageProvider({ children }: { children: ReactNode }) {
    const [language, setLanguageState] = useState<Language>('en');

    // Load language from localStorage on mount
    useEffect(() => {
        const saved = localStorage.getItem('language') as Language;
        if (saved && (saved === 'en' || saved === 'id')) {
            setLanguageState(saved);
        }
    }, []);

    const setLanguage = (lang: Language) => {
        setLanguageState(lang);
        localStorage.setItem('language', lang);
    };

    const t = (key: string): string => {
        return translations[language][key as keyof typeof translations.en] || key;
    };

    return (
        <LanguageContext.Provider value={{ language, setLanguage, t }}>
            {children}
        </LanguageContext.Provider>
    );
}

export function useLanguage() {
    const context = useContext(LanguageContext);
    if (context === undefined) {
        throw new Error('useLanguage must be used within a LanguageProvider');
    }
    return context;
}
