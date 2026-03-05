'use client';

import { createContext, useContext, useEffect, useState } from 'react';

type Theme = 'dark' | 'light';

interface ThemeContextValue {
    theme: Theme;
    toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextValue>({
    theme: 'dark',
    toggleTheme: () => { },
});

export function ThemeProvider({ children }: { children: React.ReactNode }) {
    const [theme, setTheme] = useState<Theme>('dark');

    // Init: check localStorage → system preference
    useEffect(() => {
        const saved = localStorage.getItem('theme') as Theme | null;
        if (saved === 'light' || saved === 'dark') {
            applyTheme(saved);
            setTheme(saved);
        } else {
            // Default to dark mode if no preference is set
            const initial: Theme = 'dark';
            applyTheme(initial);
            setTheme(initial);
        }
    }, []);

    function applyTheme(t: Theme) {
        if (t === 'light') {
            document.documentElement.classList.add('light');
            document.documentElement.classList.remove('dark');
        } else {
            document.documentElement.classList.remove('light');
            document.documentElement.classList.add('dark');
        }
    }

    function toggleTheme() {
        const next: Theme = theme === 'dark' ? 'light' : 'dark';
        applyTheme(next);
        setTheme(next);
        localStorage.setItem('theme', next);
    }

    return (
        <ThemeContext.Provider value={{ theme, toggleTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}

export function useTheme() {
    return useContext(ThemeContext);
}
