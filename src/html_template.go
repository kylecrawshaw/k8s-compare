package main

import (
	"fmt"
	"strconv"
	"strings"
)

// generateHTMLTemplate creates the complete HTML template with embedded data
func generateHTMLTemplate(config *ComparisonConfig, clusterAJSON, clusterBJSON, timestamp string) string {
	// Build the HTML template with proper escaping
	template := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kubernetes Resource Comparison Report - ` + timestamp + `</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f7fa; line-height: 1.6; }
        .container { max-width: 1400px; margin: 0 auto; padding: 20px; }
        .header { background: white; border-radius: 12px; padding: 30px; margin-bottom: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header h1 { color: #2c3e50; font-size: 2.5rem; margin-bottom: 10px; }
        .header p { color: #7f8c8d; font-size: 1.1rem; }
        .metadata-section { background: white; border-radius: 12px; padding: 30px; margin-bottom: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .metadata-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .metadata-card { background: #f8f9fa; border-radius: 8px; padding: 20px; border-left: 4px solid #3498db; }
        .metadata-card h3 { color: #2c3e50; margin-bottom: 15px; font-size: 1.2rem; }
        .metadata-item { margin-bottom: 10px; }
        .metadata-label { font-weight: 600; color: #34495e; margin-bottom: 5px; }
        .metadata-value { color: #7f8c8d; background: white; padding: 8px 12px; border-radius: 4px; font-family: 'Monaco', 'Consolas', monospace; font-size: 0.9rem; }
        .resource-tags { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 10px; }
        .resource-tag { background: #3498db; color: white; padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; }
        .tabs { display: flex; background: white; border-radius: 12px 12px 0 0; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .tab { flex: 1; padding: 15px 20px; background: #ecf0f1; border: none; cursor: pointer; font-weight: 600; color: #7f8c8d; transition: all 0.3s ease; }
        .tab:hover { background: #d5dbdb; }
        .tab.active { background: #3498db; color: white; }
        .tab-content { background: white; border-radius: 0 0 12px 12px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); display: none; }
        .tab-content.active { display: block; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: #f8f9fa; border-radius: 8px; padding: 20px; text-align: center; border-left: 4px solid #3498db; }
        .stat-number { font-size: 2rem; font-weight: bold; color: #2c3e50; display: block; }
        .stat-label { color: #7f8c8d; margin-top: 5px; }
        .comparison-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
        .resource-list { background: #f8f9fa; border-radius: 8px; padding: 20px; }
        .resource-list h3 { color: #2c3e50; margin-bottom: 15px; font-size: 1.2rem; }
        .resource-item { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid #ecf0f1; }
        .resource-item:last-child { border-bottom: none; }
        .resource-name { font-weight: 500; color: #2c3e50; }
        .resource-count { background: #3498db; color: white; padding: 4px 12px; border-radius: 20px; font-size: 0.9rem; font-weight: 600; }
        .resource-diff { background: #f8f9fa; border-radius: 8px; margin-bottom: 20px; overflow: hidden; border: 1px solid #ecf0f1; }
        .resource-diff.expanded .resource-content { display: block; }
        .resource-header { background: #ecf0f1; padding: 15px 20px; cursor: pointer; display: flex; justify-content: space-between; align-items: center; transition: background 0.3s ease; }
        .resource-header:hover { background: #d5dbdb; }
        .resource-content { display: none; padding: 20px; }
        .individual-resource { background: white; border-radius: 6px; margin-bottom: 15px; overflow: hidden; border: 1px solid #e0e6ed; }
        .individual-resource.expanded .individual-content { display: block; }
        .individual-header { background: #f8f9fa; padding: 12px 15px; cursor: pointer; display: flex; justify-content: space-between; align-items: center; font-size: 0.95rem; }
        .individual-header:hover { background: #e9ecef; }
        .individual-content { display: none; padding: 15px; }
        .resource-metadata { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin-bottom: 20px; padding: 15px; background: #f8f9fa; border-radius: 6px; }
        .status-badge { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; font-weight: 600; text-transform: uppercase; }
        .status-different { background: #fff3cd; color: #856404; }
        .status-only-a { background: #f8d7da; color: #721c24; }
        .status-only-b { background: #d1ecf1; color: #0c5460; }
        .diff-row { display: grid; grid-template-columns: 200px 1fr 1fr; gap: 15px; margin-bottom: 15px; padding: 15px; background: #f8f9fa; border-radius: 6px; }
        .diff-field { font-weight: 600; color: #2c3e50; align-self: start; }
        .diff-value { padding: 10px; border-radius: 4px; overflow-x: auto; max-height: 400px; overflow-y: auto; }
        .diff-value.different { background: #fff3cd; border-left: 4px solid #ffc107; }
        .diff-value.missing { background: #f8d7da; border-left: 4px solid #dc3545; }
        .diff-value.added { background: #d1ecf1; border-left: 4px solid #17a2b8; }
        .json-key { color: #0066cc; font-weight: 600; }
        .json-string { color: #008000; }
        .json-number { color: #ff6600; }
        .json-boolean { color: #cc0000; font-weight: 600; }
        .json-null { color: #999999; font-style: italic; }
        .json-object, .json-array { color: #333; }
        .loading { text-align: center; padding: 40px; color: #7f8c8d; font-size: 1.1rem; }
        .toggle-icon { transition: transform 0.3s ease; }
        .expanded .toggle-icon { transform: rotate(90deg); }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Kubernetes Resource Comparison Report</h1>
            <p>Automated comparison between ` + config.ClusterA.Context + ` and ` + config.ClusterB.Context + `</p>
            <p style="margin-top: 10px; font-size: 1rem;">Generated on ` + timestamp + `</p>
        </div>
        
        <div class="metadata-section">
            <h2 style="color: #2c3e50; margin-bottom: 20px;">📊 Comparison Metadata</h2>
            <div class="metadata-grid">
                <div class="metadata-card">
                    <h3>🅰️ Cluster A</h3>
                    <div class="metadata-item">
                        <div class="metadata-label">Context</div>
                        <div class="metadata-value">` + config.ClusterA.Context + `</div>
                    </div>
                    <div class="metadata-item">
                        <div class="metadata-label">Namespaces</div>
                        <div class="metadata-value">` + strings.Join(config.ClusterA.Namespaces, ", ") + `</div>
                    </div>
                    <div class="metadata-item">
                        <div class="metadata-label">Resource Count</div>
                        <div class="metadata-value">` + fmt.Sprintf("%d", len(config.ClusterA.Data)) + ` resources</div>
                    </div>
                    <div class="resource-tags">
                        ` + generateResourceTags(config.ClusterA.Resources) + `
                    </div>
                </div>
                
                <div class="metadata-card">
                    <h3>🅱️ Cluster B</h3>
                    <div class="metadata-item">
                        <div class="metadata-label">Context</div>
                        <div class="metadata-value">` + config.ClusterB.Context + `</div>
                    </div>
                    <div class="metadata-item">
                        <div class="metadata-label">Namespaces</div>
                        <div class="metadata-value">` + strings.Join(config.ClusterB.Namespaces, ", ") + `</div>
                    </div>
                    <div class="metadata-item">
                        <div class="metadata-label">Resource Count</div>
                        <div class="metadata-value">` + fmt.Sprintf("%d", len(config.ClusterB.Data)) + ` resources</div>
                    </div>
                    <div class="resource-tags">
                        ` + generateResourceTags(config.ClusterB.Resources) + `
                    </div>
                </div>
            </div>
        </div>

        <div style="margin-bottom: 24px; text-align: center;">
            <label style="font-weight: 500; color: #2c3e50; cursor: pointer;">
                <input type="checkbox" id="namespace-toggle" checked style="margin-right: 8px; vertical-align: middle;" />
                Use <span style="font-family: monospace;">namespace</span> for resource comparison
            </label>
        </div>

        <div class="tabs">
            <button class="tab active" onclick="showTab('overview')">📊 Overview</button>
            <button class="tab" onclick="showTab('breakdown')">📋 Resource Breakdown</button>
            <button class="tab" onclick="showTab('detailed')">🔎 Detailed Comparison</button>
        </div>

        <div id="overview" class="tab-content active">
            <div class="stats-grid" id="stats-grid"></div>
        </div>

        <div id="breakdown" class="tab-content">
            <div class="comparison-grid" id="breakdown-content"></div>
        </div>

        <div id="detailed" class="tab-content">
            <div id="detailed-content"></div>
        </div>
    </div>

    <script>
        const file1Data = ` + clusterAJSON + `;
        const file2Data = ` + clusterBJSON + `;
        let useNamespace = ` + strconv.FormatBool(config.CompareNamespaces) + `;

        document.addEventListener('DOMContentLoaded', function() {
            document.getElementById('namespace-toggle').addEventListener('change', function(e) {
                useNamespace = e.target.checked;
                performComparison();
            });
            performComparison();
        });
        
        function showTab(tabName) {
            const contents = document.querySelectorAll('.tab-content');
            contents.forEach(content => content.classList.remove('active'));
            
            const tabs = document.querySelectorAll('.tab');
            tabs.forEach(tab => tab.classList.remove('active'));
            
            document.getElementById(tabName).classList.add('active');
            event.target.classList.add('active');
        }
        
        function performComparison() {
            if (!file1Data || !file2Data) {
                console.error('Data not available for comparison');
                return;
            }
            
            const comparison = compareFiles(file1Data, file2Data);
            displayOverview(comparison);
            displayBreakdown(comparison);
            displayDetailed(comparison);
        }
        
        function compareFiles(file1, file2) {
            const file1Resources = groupByKind(file1);
            const file2Resources = groupByKind(file2);
            
            const allKinds = new Set([...Object.keys(file1Resources), ...Object.keys(file2Resources)]);
            
            const comparison = {
                totalFile1: file1.length,
                totalFile2: file2.length,
                kinds: {}
            };
            
            allKinds.forEach(kind => {
                const resources1 = file1Resources[kind] || [];
                const resources2 = file2Resources[kind] || [];
                
                comparison.kinds[kind] = {
                    file1Count: resources1.length,
                    file2Count: resources2.length,
                    differences: compareResourceLists(resources1, resources2)
                };
            });
            
            return comparison;
        }
        
        function groupByKind(resources) {
            const grouped = {};
            resources.forEach(resource => {
                const kind = resource.kind || 'Unknown';
                if (!grouped[kind]) {
                    grouped[kind] = [];
                }
                grouped[kind].push(resource);
            });
            return grouped;
        }
        
        function compareResourceLists(list1, list2) {
            const differences = [];
            const map1 = new Map();
            const map2 = new Map();
            
            list1.forEach(resource => {
                const key = getResourceKey(resource);
                map1.set(key, resource);
            });
            
            list2.forEach(resource => {
                const key = getResourceKey(resource);
                map2.set(key, resource);
            });
            
            map1.forEach((resource, key) => {
                if (!map2.has(key)) {
                    differences.push({
                        type: 'only_in_file1',
                        resource: resource,
                        name: resource.metadata?.name || 'unknown'
                    });
                }
            });
            
            map2.forEach((resource, key) => {
                if (!map1.has(key)) {
                    differences.push({
                        type: 'only_in_file2',
                        resource: resource,
                        name: resource.metadata?.name || 'unknown'
                    });
                }
            });
            
            map1.forEach((resource1, key) => {
                if (map2.has(key)) {
                    const resource2 = map2.get(key);
                    const resourceDiffs = findResourceDifferences(resource1, resource2);
                    if (resourceDiffs.length > 0) {
                        differences.push({
                            type: 'different',
                            resource1: resource1,
                            resource2: resource2,
                            name: resource1.metadata?.name || 'unknown',
                            differences: resourceDiffs
                        });
                    }
                }
            });
            
            return differences;
        }
        
        function getResourceKey(resource) {
            const name = resource.metadata?.name || 'unknown';
            const namespace = resource.metadata?.namespace || 'default';
            const kind = resource.kind || 'unknown';
            if (useNamespace) {
                return kind + '/' + namespace + '/' + name;
            } else {
                return kind + '/' + name;
            }
        }
        
        function findResourceDifferences(obj1, obj2, path = '') {
            const differences = [];
            const skipFields = ['resourceVersion', 'uid', 'generation', 'creationTimestamp', 'managedFields'];
            const allKeys = new Set([...Object.keys(obj1), ...Object.keys(obj2)]);
            
            allKeys.forEach(key => {
                if (skipFields.includes(key)) return;
                
                const newPath = path ? path + '.' + key : key;
                const val1 = obj1[key];
                const val2 = obj2[key];
                
                if (val1 === undefined && val2 !== undefined) {
                    differences.push({ field: newPath, value1: undefined, value2: val2 });
                } else if (val1 !== undefined && val2 === undefined) {
                    differences.push({ field: newPath, value1: val1, value2: undefined });
                } else if (typeof val1 === 'object' && typeof val2 === 'object' && val1 !== null && val2 !== null) {
                    differences.push(...findResourceDifferences(val1, val2, newPath));
                } else if (JSON.stringify(val1) !== JSON.stringify(val2)) {
                    differences.push({ field: newPath, value1: val1, value2: val2 });
                }
            });
            
            return differences;
        }
        
        function displayOverview(comparison) {
            const statsGrid = document.getElementById('stats-grid');
            const totalDifferences = Object.values(comparison.kinds).reduce((sum, kind) => sum + kind.differences.length, 0);
            const totalKinds = Object.keys(comparison.kinds).length;
            
            let onlyInFile1 = 0, onlyInFile2 = 0, different = 0;
            
            Object.values(comparison.kinds).forEach(kind => {
                kind.differences.forEach(diff => {
                    if (diff.type === 'only_in_file1') onlyInFile1++;
                    else if (diff.type === 'only_in_file2') onlyInFile2++;
                    else if (diff.type === 'different') different++;
                });
            });
            
            statsGrid.innerHTML = '<div class="stat-card"><span class="stat-number">' + comparison.totalFile1 + '</span><div class="stat-label">Resources in Cluster A</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + comparison.totalFile2 + '</span><div class="stat-label">Resources in Cluster B</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + totalKinds + '</span><div class="stat-label">Resource Types</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + totalDifferences + '</span><div class="stat-label">Total Differences</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + onlyInFile1 + '</span><div class="stat-label">Only in Cluster A</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + onlyInFile2 + '</span><div class="stat-label">Only in Cluster B</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + different + '</span><div class="stat-label">Different Resources</div></div>';
        }
        
        function displayBreakdown(comparison) {
            const breakdownContent = document.getElementById('breakdown-content');
            
            const file1Html = Object.entries(comparison.kinds).map(([kind, data]) => 
                '<div class="resource-item"><span class="resource-name">' + kind + '</span><span class="resource-count">' + data.file1Count + '</span></div>'
            ).join('');
            
            const file2Html = Object.entries(comparison.kinds).map(([kind, data]) => 
                '<div class="resource-item"><span class="resource-name">' + kind + '</span><span class="resource-count">' + data.file2Count + '</span></div>'
            ).join('');
            
            breakdownContent.innerHTML = '<div class="resource-list"><h3>🅰️ Cluster A Resources</h3>' + file1Html + '</div>' +
                '<div class="resource-list"><h3>🅱️ Cluster B Resources</h3>' + file2Html + '</div>';
        }
        
        function displayDetailed(comparison) {
            const detailedContent = document.getElementById('detailed-content');
            let html = '';
            
            Object.entries(comparison.kinds).forEach(([kind, data]) => {
                if (data.differences.length === 0) return;
                
                html += '<div class="resource-diff"><div class="resource-header" onclick="toggleResourceDiff(this)"><h3>' + kind + ' (' + data.differences.length + ' differences)</h3><span class="toggle-icon">▶</span></div><div class="resource-content">';
                
                data.differences.forEach(diff => {
                    let statusBadge = '';
                    if (diff.type === 'different') statusBadge = '<span class="status-badge status-different">Different</span>';
                    else if (diff.type === 'only_in_file1') statusBadge = '<span class="status-badge status-only-a">Only in A</span>';
                    else if (diff.type === 'only_in_file2') statusBadge = '<span class="status-badge status-only-b">Only in B</span>';
                    
                    html += '<div class="individual-resource"><div class="individual-header" onclick="toggleIndividualResource(this)"><span>' + diff.name + ' ' + statusBadge + '</span><span class="toggle-icon">▶</span></div><div class="individual-content">';
                    
                    if (diff.type === 'different') {
                        const resource1 = diff.resource1;
                        const resource2 = diff.resource2;
                        html += '<div class="resource-metadata"><div class="metadata-item"><div class="metadata-label">Namespace</div><div class="metadata-value">' + (resource1.metadata.namespace || 'default') + '</div></div></div>';
                        
                        diff.differences.forEach(fieldDiff => {
                            html += '<div class="diff-row"><div class="diff-field">' + fieldDiff.field + '</div><div class="diff-value different"><strong>Cluster A:</strong><br>' + renderRichJson(fieldDiff.value1) + '</div><div class="diff-value different"><strong>Cluster B:</strong><br>' + renderRichJson(fieldDiff.value2) + '</div></div>';
                        });
                    } else {
                        const resource = diff.resource;
                        const cluster = diff.type === 'only_in_file1' ? 'Cluster A' : 'Cluster B';
                        const valueClass = diff.type === 'only_in_file1' ? 'missing' : 'added';
                        html += '<div class="resource-metadata"><div class="metadata-item"><div class="metadata-label">Status</div><div class="metadata-value">Only exists in ' + cluster + '</div></div></div>';
                        html += '<div class="diff-value ' + valueClass + '"><strong>Resource Definition:</strong><br>' + renderRichJson(resource) + '</div>';
                    }
                    
                    html += '</div></div>';
                });
                
                html += '</div></div>';
            });
            
            if (html === '') {
                html = '<div class="loading">No detailed differences found</div>';
            }
            
            detailedContent.innerHTML = html;
        }
        
        function toggleResourceDiff(header) {
            header.parentElement.classList.toggle('expanded');
        }
        
        function toggleIndividualResource(header) {
            header.parentElement.classList.toggle('expanded');
        }
        
        function renderRichJson(value, depth = 0) {
            if (value === null) return '<span class="json-null">null</span>';
            if (value === undefined) return '<span class="json-null">undefined</span>';
            if (typeof value === 'string') return '<span class="json-string">"' + escapeHtml(value) + '"</span>';
            if (typeof value === 'number') return '<span class="json-number">' + value + '</span>';
            if (typeof value === 'boolean') return '<span class="json-boolean">' + value + '</span>';
            
            if (Array.isArray(value)) {
                if (value.length === 0) return '<span class="json-array">[]</span>';
                let html = '<div class="json-array">[<br>';
                value.forEach((item, index) => {
                    const indent = '&nbsp;'.repeat((depth + 1) * 2);
                    html += indent + renderRichJson(item, depth + 1);
                    if (index < value.length - 1) html += ',';
                    html += '<br>';
                });
                html += '&nbsp;'.repeat(depth * 2) + ']</div>';
                return html;
            }
            
            if (typeof value === 'object') {
                const keys = Object.keys(value);
                if (keys.length === 0) return '<span class="json-object">{}</span>';
                let html = '<div class="json-object">{<br>';
                keys.forEach((key, index) => {
                    const indent = '&nbsp;'.repeat((depth + 1) * 2);
                    html += indent + '<span class="json-key">"' + escapeHtml(key) + '"</span>: ' + renderRichJson(value[key], depth + 1);
                    if (index < keys.length - 1) html += ',';
                    html += '<br>';
                });
                html += '&nbsp;'.repeat(depth * 2) + '}</div>';
                return html;
            }
            
            return escapeHtml(String(value));
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
    </script>
</body>
</html>`

	return template
}

// generateJavaScriptFunctions returns the JavaScript functions for the HTML report
func generateJavaScriptFunctions() string {
	return `
        
        function showTab(tabName) {
            const contents = document.querySelectorAll('.tab-content');
            contents.forEach(content => content.classList.remove('active'));
            
            const tabs = document.querySelectorAll('.tab');
            tabs.forEach(tab => tab.classList.remove('active'));
            
            document.getElementById(tabName).classList.add('active');
            event.target.classList.add('active');
        }
        
        function performComparison() {
            if (!file1Data || !file2Data) {
                console.error('Data not available for comparison');
                return;
            }
            
            const comparison = compareFiles(file1Data, file2Data);
            displayOverview(comparison);
            displayBreakdown(comparison);
            displayDetailed(comparison);
        }
        
        function compareFiles(file1, file2) {
            const file1Resources = groupByKind(file1);
            const file2Resources = groupByKind(file2);
            
            const allKinds = new Set([...Object.keys(file1Resources), ...Object.keys(file2Resources)]);
            
            const comparison = {
                totalFile1: file1.length,
                totalFile2: file2.length,
                kinds: {}
            };
            
            allKinds.forEach(kind => {
                const resources1 = file1Resources[kind] || [];
                const resources2 = file2Resources[kind] || [];
                
                comparison.kinds[kind] = {
                    file1Count: resources1.length,
                    file2Count: resources2.length,
                    differences: compareResourceLists(resources1, resources2)
                };
            });
            
            return comparison;
        }
        
        function groupByKind(resources) {
            const grouped = {};
            resources.forEach(resource => {
                const kind = resource.kind || 'Unknown';
                if (!grouped[kind]) {
                    grouped[kind] = [];
                }
                grouped[kind].push(resource);
            });
            return grouped;
        }
        
        function compareResourceLists(list1, list2) {
            const differences = [];
            const map1 = new Map();
            const map2 = new Map();
            
            list1.forEach(resource => {
                const key = getResourceKey(resource);
                map1.set(key, resource);
            });
            
            list2.forEach(resource => {
                const key = getResourceKey(resource);
                map2.set(key, resource);
            });
            
            map1.forEach((resource, key) => {
                if (!map2.has(key)) {
                    differences.push({
                        type: 'only_in_file1',
                        resource: resource,
                        name: resource.metadata?.name || 'unknown'
                    });
                }
            });
            
            map2.forEach((resource, key) => {
                if (!map1.has(key)) {
                    differences.push({
                        type: 'only_in_file2',
                        resource: resource,
                        name: resource.metadata?.name || 'unknown'
                    });
                }
            });
            
            map1.forEach((resource1, key) => {
                if (map2.has(key)) {
                    const resource2 = map2.get(key);
                    const resourceDiffs = findResourceDifferences(resource1, resource2);
                    if (resourceDiffs.length > 0) {
                        differences.push({
                            type: 'different',
                            resource1: resource1,
                            resource2: resource2,
                            name: resource1.metadata?.name || 'unknown',
                            differences: resourceDiffs
                        });
                    }
                }
            });
            
            return differences;
        }
        
        function getResourceKey(resource) {
            const name = resource.metadata?.name || 'unknown';
            const namespace = resource.metadata?.namespace || 'default';
            const kind = resource.kind || 'unknown';
            if (useNamespace) {
                return kind + '/' + namespace + '/' + name;
            } else {
                return kind + '/' + name;
            }
        }
        
        function findResourceDifferences(obj1, obj2, path = '') {
            const differences = [];
            const skipFields = ['resourceVersion', 'uid', 'generation', 'creationTimestamp', 'managedFields'];
            const allKeys = new Set([...Object.keys(obj1), ...Object.keys(obj2)]);
            
            allKeys.forEach(key => {
                if (skipFields.includes(key)) return;
                
                const newPath = path ? path + '.' + key : key;
                const val1 = obj1[key];
                const val2 = obj2[key];
                
                if (val1 === undefined && val2 !== undefined) {
                    differences.push({ field: newPath, value1: undefined, value2: val2 });
                } else if (val1 !== undefined && val2 === undefined) {
                    differences.push({ field: newPath, value1: val1, value2: undefined });
                } else if (typeof val1 === 'object' && typeof val2 === 'object' && val1 !== null && val2 !== null) {
                    differences.push(...findResourceDifferences(val1, val2, newPath));
                } else if (JSON.stringify(val1) !== JSON.stringify(val2)) {
                    differences.push({ field: newPath, value1: val1, value2: val2 });
                }
            });
            
            return differences;
        }
        
        function displayOverview(comparison) {
            const statsGrid = document.getElementById('stats-grid');
            const totalDifferences = Object.values(comparison.kinds).reduce((sum, kind) => sum + kind.differences.length, 0);
            const totalKinds = Object.keys(comparison.kinds).length;
            
            let onlyInFile1 = 0, onlyInFile2 = 0, different = 0;
            
            Object.values(comparison.kinds).forEach(kind => {
                kind.differences.forEach(diff => {
                    if (diff.type === 'only_in_file1') onlyInFile1++;
                    else if (diff.type === 'only_in_file2') onlyInFile2++;
                    else if (diff.type === 'different') different++;
                });
            });
            
            statsGrid.innerHTML = '<div class="stat-card"><span class="stat-number">' + comparison.totalFile1 + '</span><div class="stat-label">Resources in Cluster A</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + comparison.totalFile2 + '</span><div class="stat-label">Resources in Cluster B</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + totalKinds + '</span><div class="stat-label">Resource Types</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + totalDifferences + '</span><div class="stat-label">Total Differences</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + onlyInFile1 + '</span><div class="stat-label">Only in Cluster A</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + onlyInFile2 + '</span><div class="stat-label">Only in Cluster B</div></div>' +
                '<div class="stat-card"><span class="stat-number">' + different + '</span><div class="stat-label">Different Resources</div></div>';
        }
        
        function displayBreakdown(comparison) {
            const breakdownContent = document.getElementById('breakdown-content');
            
            const file1Html = Object.entries(comparison.kinds).map(([kind, data]) => 
                '<div class="resource-item"><span class="resource-name">' + kind + '</span><span class="resource-count">' + data.file1Count + '</span></div>'
            ).join('');
            
            const file2Html = Object.entries(comparison.kinds).map(([kind, data]) => 
                '<div class="resource-item"><span class="resource-name">' + kind + '</span><span class="resource-count">' + data.file2Count + '</span></div>'
            ).join('');
            
            breakdownContent.innerHTML = '<div class="resource-list"><h3>🅰️ Cluster A Resources</h3>' + file1Html + '</div>' +
                '<div class="resource-list"><h3>🅱️ Cluster B Resources</h3>' + file2Html + '</div>';
        }
        
        function displayDetailed(comparison) {
            const detailedContent = document.getElementById('detailed-content');
            let html = '';
            
            Object.entries(comparison.kinds).forEach(([kind, data]) => {
                if (data.differences.length === 0) return;
                
                html += '<div class="resource-diff"><div class="resource-header" onclick="toggleResourceDiff(this)"><h3>' + kind + ' (' + data.differences.length + ' differences)</h3><span class="toggle-icon">▶</span></div><div class="resource-content">';
                
                data.differences.forEach(diff => {
                    let statusBadge = '';
                    if (diff.type === 'different') statusBadge = '<span class="status-badge status-different">Different</span>';
                    else if (diff.type === 'only_in_file1') statusBadge = '<span class="status-badge status-only-a">Only in A</span>';
                    else if (diff.type === 'only_in_file2') statusBadge = '<span class="status-badge status-only-b">Only in B</span>';
                    
                    html += '<div class="individual-resource"><div class="individual-header" onclick="toggleIndividualResource(this)"><span>' + diff.name + ' ' + statusBadge + '</span><span class="toggle-icon">▶</span></div><div class="individual-content">';
                    
                    if (diff.type === 'different') {
                        const resource1 = diff.resource1;
                        const resource2 = diff.resource2;
                        html += '<div class="resource-metadata"><div class="metadata-item"><div class="metadata-label">Namespace</div><div class="metadata-value">' + (resource1.metadata.namespace || 'default') + '</div></div></div>';
                        
                        diff.differences.forEach(fieldDiff => {
                            html += '<div class="diff-row"><div class="diff-field">' + fieldDiff.field + '</div><div class="diff-value different"><strong>Cluster A:</strong><br>' + renderRichJson(fieldDiff.value1) + '</div><div class="diff-value different"><strong>Cluster B:</strong><br>' + renderRichJson(fieldDiff.value2) + '</div></div>';
                        });
                    } else {
                        const resource = diff.resource;
                        const cluster = diff.type === 'only_in_file1' ? 'Cluster A' : 'Cluster B';
                        const valueClass = diff.type === 'only_in_file1' ? 'missing' : 'added';
                        html += '<div class="resource-metadata"><div class="metadata-item"><div class="metadata-label">Status</div><div class="metadata-value">Only exists in ' + cluster + '</div></div></div>';
                        html += '<div class="diff-value ' + valueClass + '"><strong>Resource Definition:</strong><br>' + renderRichJson(resource) + '</div>';
                    }
                    
                    html += '</div></div>';
                });
                
                html += '</div></div>';
            });
            
            if (html === '') {
                html = '<div class="loading">No detailed differences found</div>';
            }
            
            detailedContent.innerHTML = html;
        }
        
        function toggleResourceDiff(header) {
            header.parentElement.classList.toggle('expanded');
        }
        
        function toggleIndividualResource(header) {
            header.parentElement.classList.toggle('expanded');
        }
        
        function renderRichJson(value, depth = 0) {
            if (value === null) return '<span class="json-null">null</span>';
            if (value === undefined) return '<span class="json-null">undefined</span>';
            if (typeof value === 'string') return '<span class="json-string">"' + escapeHtml(value) + '"</span>';
            if (typeof value === 'number') return '<span class="json-number">' + value + '</span>';
            if (typeof value === 'boolean') return '<span class="json-boolean">' + value + '</span>';
            
            if (Array.isArray(value)) {
                if (value.length === 0) return '<span class="json-array">[]</span>';
                let html = '<div class="json-array">[<br>';
                value.forEach((item, index) => {
                    const indent = '&nbsp;'.repeat((depth + 1) * 2);
                    html += indent + renderRichJson(item, depth + 1);
                    if (index < value.length - 1) html += ',';
                    html += '<br>';
                });
                html += '&nbsp;'.repeat(depth * 2) + ']</div>';
                return html;
            }
            
            if (typeof value === 'object') {
                const keys = Object.keys(value);
                if (keys.length === 0) return '<span class="json-object">{}</span>';
                let html = '<div class="json-object">{<br>';
                keys.forEach((key, index) => {
                    const indent = '&nbsp;'.repeat((depth + 1) * 2);
                    html += indent + '<span class="json-key">"' + escapeHtml(key) + '"</span>: ' + renderRichJson(value[key], depth + 1);
                    if (index < keys.length - 1) html += ',';
                    html += '<br>';
                });
                html += '&nbsp;'.repeat(depth * 2) + '}</div>';
                return html;
            }
            
            return escapeHtml(String(value));
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }`
}
