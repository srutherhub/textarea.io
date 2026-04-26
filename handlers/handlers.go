package handlers

import (
	"app/components"
	"app/services"
	"fmt"
	"net/http"

	v "github.com/srutherhub/web-app/views"
)

func Base() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			indexProps := v.IndexProps{Title: "Oops! Something went wrong."}
			page := components.Error()
			v.Index(page, indexProps).Render(r.Context(), w)
			return
		}
	}
}

func GetCreateSpace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexProps := v.IndexProps{Title: "textarea.io"}
		page := components.CreateSpace()
		v.Index(page, indexProps).Render(r.Context(), w)
	}
}

func CreateSpace(as *services.AppService, au *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		id, key, name, err := as.CreateSpace(name)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			components.CreateSpaceFail().Render(r.Context(), w)
			return
		}

		host := r.Host
		url := host + "/space/" + id

		sessionCookie, err := r.Cookie("session")

		if err != nil {
			token, _ := au.CreateSessionToken(services.SpaceInfo{Name: name, ID: id, Key: key})
			SetSessionCookie(w, token)
		} else {
			token, _ := au.AddToSessionToken(sessionCookie.Value, services.SpaceInfo{Name: name, ID: id, Key: key})
			SetSessionCookie(w, token)
		}

		components.CreateSpaceSuccess(id, name, key, url).Render(r.Context(), w)
	}
}

func GetSpace(as *services.AppService, au *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spaceId := r.PathValue("id")
		space, err := as.GetByID(spaceId)

		if err != nil {
			indexProps := v.IndexProps{Title: "Oops! Something went wrong."}
			page := components.Error()
			v.Index(page, indexProps).Render(r.Context(), w)
			return
		}

		indexProps := v.IndexProps{Title: space.Name}
		page := components.Space(components.ISpaceProps{ID: space.ID, Name: space.Name, Content: space.Content, Code: space.Key})
		v.Index(page, indexProps).Render(r.Context(), w)
	}
}

func SaveSpace(as *services.AppService, au *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spaceId := r.PathValue("id")
		content := r.FormValue("content")

		err := as.AddContent(spaceId, content)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to save")
			return
		}
	}
}

func DeleteSpace(as *services.AppService, au *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spaceId := r.PathValue("id")
		content := r.FormValue("content")

		err := as.AddContent(spaceId, content)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to save")
			return
		}
	}
}

func Area(au *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionToken := GetSessionTokenFromCookie(au, r)

		indexProps := v.IndexProps{Title: "Area"}
		page := components.Area(components.AreaProps{Areas: sessionToken})
		v.Index(page, indexProps).Render(r.Context(), w)
	}
}

func AreaAuth(au *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		key := r.FormValue("key")

		areaInfo, err := au.Authenticate(id, key)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized:"+err.Error())
			return
		}
		AddOrCreateSessionCookie(areaInfo.Name, areaInfo.ID, areaInfo.Key, au, w, r)
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/area/"+areaInfo.ID)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Redirect(w, r, "/area/"+areaInfo.ID, http.StatusFound)
	}
}

func GetSessionTokenFromCookie(au *services.AuthService, r *http.Request) services.SessionToken {
	session, err := r.Cookie("session")
	if err != nil || session.Value == "" {
		return services.SessionToken{}
	}

	token, err := au.DecryptTokens(session.Value)
	if err != nil {
		return services.SessionToken{}
	}

	return token
}

func AddOrCreateSessionCookie(name string, id string, key string, au *services.AuthService, w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")

	if err != nil {
		token, _ := au.CreateSessionToken(services.SpaceInfo{Name: name, ID: id, Key: key})
		SetSessionCookie(w, token)
	} else {
		token, _ := au.AddToSessionToken(sessionCookie.Value, services.SpaceInfo{Name: name, ID: id, Key: key})
		SetSessionCookie(w, token)
	}
}

func SetSessionCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}
