package library

import (
    "errors"
)

type MusicEntry struct {
    Id     string
    Name   string
    Artist string
    Source string
    Type   string
}

type MusicManager struct {
    musics []MusicEntry
}

func NewMusicManager() *MusicManager {
    return &MusicManager{make([]MusicEntry, 0)}
}

func (m *MusicManager) Len() int {
    return len(m.musics)
}

func (m *MusicManager) Get(index int) (music *MusicEntry, err error) {
    if index < 0 || index >= len(m.musics) {
        return nil, errors.New("Index out of range.")
    }
    return &m.musics[index], nil
}

func (m *MusicManager) Find(name string) *MusicEntry {
    if len(m.musics) == 0 {
        return nil
    }

    for _, m := range m.musics {
        if m.Name == name {
            return &m
        }
    }
    return nil
}

func (m *MusicManager) Add(music *MusicEntry) {
    m.musics = append(m.musics, *music)
}

func (m *MusicManager) Remove(name string) *MusicEntry {
    index := -1
    for i := 0; i < len(m.musics); i++ {
        if m.musics[i].Name == name {
            index = i
        }
    }

    if index == -1 {
        return nil
    }

    removedMusic := &m.musics[index]

    // 从数组切片中删除元素
    if len(m.musics) == 1 { // 删除仅有的一个元素
        m.musics = make([]MusicEntry, 0)
    } else {
        if index == 0 { // 删除第一个元素
            m.musics = m.musics[1:]
        } else if index == len(m.musics)-1 { // 删除最后一个元素
            m.musics = m.musics[:index-1]
        } else { // 删除中间元素
            m.musics = append(m.musics[:index], m.musics[index+1:]...)
        }
    }

    return removedMusic
}
