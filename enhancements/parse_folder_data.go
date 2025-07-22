package aep

import "github.com/rioam2/rifx"

// Enhancement: Parse folder data blocks (fdta)
// Extends existing folder parsing capabilities

// FolderMetadata represents additional folder properties
type FolderMetadata struct {
	IsExpanded bool
	Color      byte // Label color index
	IsShyred   bool // Hidden in timeline
}

// parseFolderData extends folder parsing following established patterns
// Reuses binary struct parsing approach from CDTA/IDTA blocks
func parseFolderData(block *rifx.Block) (*FolderMetadata, error) {
	// Following the struct parsing pattern from lines 136-155 in item.go
	type FDTA struct {
		Flags    byte     // Bit flags for folder state
		Color    byte     // Label color (0-15)
		Reserved [14]byte // Padding/future use
	}
	
	fdtaDesc := &FDTA{}
	err := block.ToStruct(fdtaDesc)
	if err != nil {
		return nil, err
	}
	
	// Extract metadata following bit flag patterns
	metadata := &FolderMetadata{
		IsExpanded: (fdtaDesc.Flags & 0x01) != 0,
		IsShyred:   (fdtaDesc.Flags & 0x02) != 0,
		Color:      fdtaDesc.Color,
	}
	
	return metadata, nil
}

// Integration point: Add to parseItem at the folder case (line 89)
// This extends the existing folder parsing logic
func enhanceFolderParsing(itemHead *rifx.List, item *Item) error {
	// Use existing FindByType pattern
	fdtaBlock, err := itemHead.FindByType("fdta")
	if err == nil { // Optional block, don't fail if missing
		metadata, err := parseFolderData(fdtaBlock)
		if err != nil {
			return err
		}
		// Would add to item.FolderMetadata field (requires Item struct extension)
		_ = metadata // Placeholder for integration
	}
	
	return nil
}