/*
   file:           task/workcenter_resolution.go
   description:    Helper untuk menentukan workcenter yang digunakan mesin
   created:        penyesuaian logika ecbtest
*/

package task

import (
	"database/sql"
	"net"
	"strings"

	"go-ecb/configs"
)

// resolveWorkcenters menentukan daftar workcenter berdasarkan IP lokal atau konfigurasi.
func resolveWorkcenters(local *sql.DB, simoCfg configs.SimoConfig) []string {
	// coba gunakan workcenter yang terdaftar pada ecbstations berdasar IP
	if local != nil {
		if ip := firstIPv4(); ip != "" {
			var raw sql.NullString
			if err := local.QueryRow("SELECT workcenters FROM ecbstations WHERE ipaddress = ? LIMIT 1", ip).Scan(&raw); err == nil && raw.Valid {
				if wcs := splitAndTrim(raw.String); len(wcs) > 0 {
					return wcs
				}
			}
		}
	}

	// fallback ke konfigurasi ENV
	if wc := strings.TrimSpace(simoCfg.EcbWorkcenters); wc != "" {
		if wcs := splitAndTrim(wc); len(wcs) > 0 {
			return wcs
		}
	}

	// default bila tidak ada konfigurasi
	return []string{"Pxxxxxxx"}
}

func splitAndTrim(csv string) []string {
	raw := strings.Split(csv, ",")
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func firstIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}
	return ""
}
